// contractgen 根据 contract.Backend 的 //contract:route 指令生成 direct 分发、HTTP 路由与前端 TS。
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const directivePrefix = "//contract:route "

type route struct {
	method string
	path   string
	status int
	auth   string
}

type method struct {
	name   string
	doc    []string
	route  route
	hasOut bool
	outTyp string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "contractgen:", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	contractDir := filepath.Join(root, "internal", "contract")
	methods, err := parseBackend(contractDir)
	if err != nil {
		return err
	}
	goFiles := map[string][]byte{
		filepath.Join(root, "internal", "direct", "backend_gen.go"): generateDirect(methods),
		filepath.Join(root, "internal", "api", "routes_gen.go"):     generateHTTP(methods),
	}
	for path, source := range goFiles {
		formatted, err := format.Source(source)
		if err != nil {
			return fmt.Errorf("format %s: %w\n%s", path, err, source)
		}
		if err := writeIfChanged(path, formatted); err != nil {
			return err
		}
	}
	generated := filepath.Join(root, "frontend", "src", "api", "generated")
	if err := os.MkdirAll(generated, 0o755); err != nil {
		return err
	}
	if err := writeIfChanged(filepath.Join(generated, "contract.ts"), generateTSContract(methods)); err != nil {
		return err
	}
	return writeIfChanged(filepath.Join(generated, "operations.ts"), generateTSOperations(methods))
}

func parseBackend(dir string) ([]method, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var files []*ast.File
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			files = append(files, f)
		}
	}
	var iface *ast.InterfaceType
	var ifaceDoc map[*ast.Field][]string
	for _, f := range files {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != "Backend" {
					continue
				}
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("Backend is not an interface")
				}
				iface = it
				ifaceDoc = methodDocs(f, it)
			}
		}
	}
	if iface == nil {
		return nil, fmt.Errorf("Backend interface not found in %s", dir)
	}
	var methods []method
	for _, field := range iface.Methods.List {
		fn, ok := field.Type.(*ast.FuncType)
		if !ok || len(field.Names) != 1 {
			continue
		}
		name := field.Names[0].Name
		docs := ifaceDoc[field]
		rt, err := parseRoute(docs)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hasOut, outTyp, err := parseResult(fn)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		methods = append(methods, method{name: name, doc: docs, route: rt, hasOut: hasOut, outTyp: outTyp})
	}
	if len(methods) == 0 {
		return nil, fmt.Errorf("Backend has no methods")
	}
	return methods, nil
}

func methodDocs(f *ast.File, iface *ast.InterfaceType) map[*ast.Field][]string {
	out := map[*ast.Field][]string{}
	for _, field := range iface.Methods.List {
		var lines []string
		if field.Doc != nil {
			for _, c := range field.Doc.List {
				lines = append(lines, strings.TrimSpace(strings.TrimPrefix(c.Text, "//")))
			}
		}
		out[field] = lines
	}
	_ = f
	return out
}

func parseRoute(docs []string) (route, error) {
	var line string
	for _, d := range docs {
		if strings.HasPrefix("//"+d, directivePrefix) || strings.HasPrefix(d, "contract:route ") {
			line = d
			break
		}
		if strings.HasPrefix(d, "contract:route ") {
			line = d
			break
		}
	}
	// docs stored without //
	for _, d := range docs {
		if strings.HasPrefix(d, "contract:route ") {
			line = d
			break
		}
	}
	if line == "" {
		return route{}, fmt.Errorf("missing //contract:route directive")
	}
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return route{}, fmt.Errorf("invalid route directive %q", line)
	}
	rt := route{method: fields[1], path: fields[2], status: 200, auth: "user"}
	for _, opt := range fields[3:] {
		switch {
		case strings.HasPrefix(opt, "status="):
			n, err := strconv.Atoi(strings.TrimPrefix(opt, "status="))
			if err != nil {
				return route{}, fmt.Errorf("invalid status: %w", err)
			}
			rt.status = n
		case strings.HasPrefix(opt, "auth="):
			rt.auth = strings.TrimPrefix(opt, "auth=")
			switch rt.auth {
			case "public", "user", "admin":
			default:
				return route{}, fmt.Errorf("invalid auth %q", rt.auth)
			}
		default:
			return route{}, fmt.Errorf("unknown route option %q", opt)
		}
	}
	return rt, nil
}

func parseResult(fn *ast.FuncType) (bool, string, error) {
	if fn.Results == nil || len(fn.Results.List) == 0 {
		return false, "", fmt.Errorf("method must return error")
	}
	results := fn.Results.List
	if len(results) == 1 {
		if ident, ok := results[0].Type.(*ast.Ident); ok && ident.Name == "error" {
			return false, "", nil
		}
		return false, "", fmt.Errorf("single result must be error")
	}
	if len(results) != 2 {
		return false, "", fmt.Errorf("want (T, error) or error")
	}
	if ident, ok := results[1].Type.(*ast.Ident); !ok || ident.Name != "error" {
		return false, "", fmt.Errorf("last result must be error")
	}
	switch t := results[0].Type.(type) {
	case *ast.Ident:
		return true, t.Name, nil
	default:
		return false, "", fmt.Errorf("unsupported result type")
	}
}

func generateDirect(methods []method) []byte {
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	b.WriteString("package direct\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n\n")
	b.WriteString("\t\"github.com/runforyou-ai/lite-ai-gateway/internal/contract\"\n")
	b.WriteString(")\n\n")
	b.WriteString("// Backend 由契约生成的认证分发层，嵌入手写实现。\n")
	b.WriteString("type Backend struct {\n\tImpl contract.Backend\n}\n\n")
	for _, m := range methods {
		writeMethodDoc(&b, m)
		if m.hasOut {
			fmt.Fprintf(&b, "func (b *Backend) %s(ctx context.Context, meta contract.RequestMeta) (contract.%s, error) {\n", m.name, m.outTyp)
			fmt.Fprintf(&b, "\tif err := authorize(ctx, meta, %q); err != nil {\n\t\tvar zero contract.%s\n\t\treturn zero, err\n\t}\n", m.route.auth, m.outTyp)
			fmt.Fprintf(&b, "\treturn b.Impl.%s(ctx, meta)\n}\n\n", m.name)
		} else {
			fmt.Fprintf(&b, "func (b *Backend) %s(ctx context.Context, meta contract.RequestMeta) error {\n", m.name)
			fmt.Fprintf(&b, "\tif err := authorize(ctx, meta, %q); err != nil {\n\t\treturn err\n\t}\n", m.route.auth)
			fmt.Fprintf(&b, "\treturn b.Impl.%s(ctx, meta)\n}\n\n", m.name)
		}
	}
	return []byte(b.String())
}

func generateHTTP(methods []method) []byte {
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	b.WriteString("package api\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"net/http\"\n\n")
	b.WriteString("\t\"github.com/runforyou-ai/lite-ai-gateway/internal/direct\"\n")
	b.WriteString(")\n\n")
	b.WriteString("// registerGeneratedRoutes 注册契约生成的 /api 路由。\n")
	b.WriteString("func registerGeneratedRoutes(mux *http.ServeMux, backend *direct.Backend) {\n")
	for _, m := range methods {
		pattern := m.route.method + " /api" + m.route.path
		fmt.Fprintf(&b, "\tmux.HandleFunc(%q, func(w http.ResponseWriter, r *http.Request) {\n", pattern)
		b.WriteString("\t\tmeta := requestMeta(r)\n")
		if m.hasOut {
			fmt.Fprintf(&b, "\t\tout, err := backend.%s(r.Context(), meta)\n", m.name)
			fmt.Fprintf(&b, "\t\twriteJSON(w, %d, out, err)\n", m.route.status)
		} else {
			fmt.Fprintf(&b, "\t\terr := backend.%s(r.Context(), meta)\n", m.name)
			fmt.Fprintf(&b, "\t\twriteJSON(w, %d, nil, err)\n", m.route.status)
		}
		b.WriteString("\t})\n")
	}
	b.WriteString("}\n")
	return []byte(b.String())
}

func generateTSContract(methods []method) []byte {
	types := map[string]bool{}
	for _, m := range methods {
		if m.hasOut {
			types[m.outTyp] = true
		}
	}
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	if types["HealthStatus"] {
		b.WriteString("export type HealthStatus = {\n  ok: boolean\n}\n\n")
	}
	if types["InstallationStatus"] {
		b.WriteString("export type InstallationStatus = {\n  installed: boolean\n  siteName: string\n}\n\n")
	}
	_ = methods
	return []byte(b.String())
}

func generateTSOperations(methods []method) []byte {
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	b.WriteString("import { call } from \"@/api/client\"\n")
	b.WriteString("import type * as contract from \"./contract\"\n\n")
	for _, m := range methods {
		fn := lowerFirst(m.name)
		for _, line := range docLines(m) {
			fmt.Fprintf(&b, "/** %s */\n", line)
		}
		if m.hasOut {
			fmt.Fprintf(&b, "export function %s(signal?: AbortSignal) {\n", fn)
			fmt.Fprintf(&b, "  return call<contract.%s>({ method: %q, path: %q }, signal)\n}\n\n", m.outTyp, m.route.method, "/api"+m.route.path)
		} else {
			fmt.Fprintf(&b, "export function %s(signal?: AbortSignal) {\n", fn)
			fmt.Fprintf(&b, "  return call<void>({ method: %q, path: %q }, signal)\n}\n\n", m.route.method, "/api"+m.route.path)
		}
	}
	return []byte(b.String())
}

func writeMethodDoc(b *strings.Builder, m method) {
	for i, line := range docLines(m) {
		if i == 0 {
			fmt.Fprintf(b, "// %s %s\n", m.name, line)
		} else {
			fmt.Fprintf(b, "// %s\n", line)
		}
	}
}

func docLines(m method) []string {
	var lines []string
	for _, d := range m.doc {
		if strings.HasPrefix(d, "contract:route ") {
			continue
		}
		lines = append(lines, d)
	}
	if len(lines) == 0 {
		return []string{m.name}
	}
	// strip leading method name if present
	first := lines[0]
	if strings.HasPrefix(first, m.name+" ") {
		lines[0] = strings.TrimPrefix(first, m.name+" ")
	}
	return lines
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}

func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func writeIfChanged(path string, content []byte) error {
	if existing, err := os.ReadFile(path); err == nil && bytes.Equal(existing, content) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}
