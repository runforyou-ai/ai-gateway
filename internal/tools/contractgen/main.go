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
	hasIn  bool
	inTyp  string
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
	mod, err := modulePath(root)
	if err != nil {
		return err
	}
	contractDir := filepath.Join(root, "internal", "contract")
	methods, err := parseBackend(contractDir)
	if err != nil {
		return err
	}
	structs, err := parseJSONStructs(contractDir)
	if err != nil {
		return err
	}
	goFiles := map[string][]byte{
		filepath.Join(root, "internal", "direct", "backend_gen.go"): generateDirect(mod, methods),
		filepath.Join(root, "internal", "api", "routes_gen.go"):     generateHTTP(mod, methods),
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
	if err := writeIfChanged(filepath.Join(generated, "contract.ts"), generateTSContract(methods, structs)); err != nil {
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
				ifaceDoc = methodDocs(it)
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
		hasIn, inTyp, err := parseInput(fn)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hasOut, outTyp, err := parseResult(fn)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		methods = append(methods, method{
			name: name, doc: docs, route: rt,
			hasIn: hasIn, inTyp: inTyp, hasOut: hasOut, outTyp: outTyp,
		})
	}
	if len(methods) == 0 {
		return nil, fmt.Errorf("Backend has no methods")
	}
	return methods, nil
}

type jsonField struct {
	name string
	typ  string
}

func parseJSONStructs(dir string) (map[string][]jsonField, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	out := map[string][]jsonField{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.TYPE {
					continue
				}
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}
					var fields []jsonField
					for _, field := range st.Fields.List {
						if len(field.Names) == 0 {
							continue
						}
						jsonName := field.Names[0].Name
						if field.Tag != nil {
							tag := strings.Trim(field.Tag.Value, "`")
							for _, part := range strings.Fields(tag) {
								if strings.HasPrefix(part, "json:\"") {
									v := strings.TrimPrefix(part, "json:\"")
									v = strings.TrimSuffix(v, "\"")
									if i := strings.IndexByte(v, ','); i >= 0 {
										v = v[:i]
									}
									if v != "" && v != "-" {
										jsonName = v
									}
								}
							}
						}
						tsType := "unknown"
						if ident, ok := field.Type.(*ast.Ident); ok {
							switch ident.Name {
							case "string":
								tsType = "string"
							case "bool":
								tsType = "boolean"
							case "int", "int64", "int32":
								tsType = "number"
							}
						}
						fields = append(fields, jsonField{name: jsonName, typ: tsType})
					}
					out[ts.Name.Name] = fields
				}
			}
		}
	}
	return out, nil
}

func methodDocs(iface *ast.InterfaceType) map[*ast.Field][]string {
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
	return out
}

func parseRoute(docs []string) (route, error) {
	var line string
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

func parseInput(fn *ast.FuncType) (bool, string, error) {
	if fn.Params == nil || len(fn.Params.List) < 2 {
		return false, "", fmt.Errorf("want at least (context.Context, RequestMeta)")
	}
	params := fn.Params.List
	if len(params) == 2 {
		return false, "", nil
	}
	if len(params) != 3 {
		return false, "", fmt.Errorf("want (ctx, meta) or (ctx, meta, Input)")
	}
	switch t := params[2].Type.(type) {
	case *ast.Ident:
		return true, t.Name, nil
	default:
		return false, "", fmt.Errorf("unsupported input type")
	}
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

func generateDirect(mod string, methods []method) []byte {
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	b.WriteString("package direct\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n\n")
	fmt.Fprintf(&b, "\t%q\n", mod+"/internal/contract")
	b.WriteString(")\n\n")
	b.WriteString("// Backend 由契约生成的认证分发层，嵌入手写实现。\n")
	b.WriteString("type Backend struct {\n\tImpl contract.Backend\n}\n\n")
	for _, m := range methods {
		writeMethodDoc(&b, m)
		sigIn := ""
		callIn := ""
		if m.hasIn {
			sigIn = ", input contract." + m.inTyp
			callIn = ", input"
		}
		if m.hasOut {
			fmt.Fprintf(&b, "func (b *Backend) %s(ctx context.Context, meta contract.RequestMeta%s) (contract.%s, error) {\n", m.name, sigIn, m.outTyp)
			fmt.Fprintf(&b, "\tif err := authorize(ctx, meta, %q); err != nil {\n\t\tvar zero contract.%s\n\t\treturn zero, err\n\t}\n", m.route.auth, m.outTyp)
			fmt.Fprintf(&b, "\treturn b.Impl.%s(ctx, meta%s)\n}\n\n", m.name, callIn)
		} else {
			fmt.Fprintf(&b, "func (b *Backend) %s(ctx context.Context, meta contract.RequestMeta%s) error {\n", m.name, sigIn)
			fmt.Fprintf(&b, "\tif err := authorize(ctx, meta, %q); err != nil {\n\t\treturn err\n\t}\n", m.route.auth)
			fmt.Fprintf(&b, "\treturn b.Impl.%s(ctx, meta%s)\n}\n\n", m.name, callIn)
		}
	}
	return []byte(b.String())
}

func generateHTTP(mod string, methods []method) []byte {
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	b.WriteString("package api\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"net/http\"\n\n")
	fmt.Fprintf(&b, "\t%q\n", mod+"/internal/contract")
	fmt.Fprintf(&b, "\t%q\n", mod+"/internal/direct")
	b.WriteString(")\n\n")
	b.WriteString("// registerGeneratedRoutes 注册契约生成的 /api 路由。\n")
	b.WriteString("func registerGeneratedRoutes(mux *http.ServeMux, backend *direct.Backend) {\n")
	for _, m := range methods {
		pattern := m.route.method + " /api" + m.route.path
		fmt.Fprintf(&b, "\tmux.HandleFunc(%q, func(w http.ResponseWriter, r *http.Request) {\n", pattern)
		b.WriteString("\t\tmeta := requestMeta(r)\n")
		callArgs := "r.Context(), meta"
		if m.hasIn {
			fmt.Fprintf(&b, "\t\tvar input contract.%s\n", m.inTyp)
			b.WriteString("\t\tif err := decodeJSONBody(r, &input); err != nil {\n")
			b.WriteString("\t\t\twriteJSON(w, http.StatusBadRequest, nil, err)\n")
			b.WriteString("\t\t\treturn\n")
			b.WriteString("\t\t}\n")
			callArgs += ", input"
		}
		if m.hasOut {
			fmt.Fprintf(&b, "\t\tout, err := backend.%s(%s)\n", m.name, callArgs)
			fmt.Fprintf(&b, "\t\twriteJSON(w, %d, out, err)\n", m.route.status)
		} else {
			fmt.Fprintf(&b, "\t\terr := backend.%s(%s)\n", m.name, callArgs)
			fmt.Fprintf(&b, "\t\twriteJSON(w, %d, nil, err)\n", m.route.status)
		}
		b.WriteString("\t})\n")
	}
	b.WriteString("}\n")
	return []byte(b.String())
}

func generateTSContract(methods []method, structs map[string][]jsonField) []byte {
	needed := map[string]bool{}
	for _, m := range methods {
		if m.hasIn {
			needed[m.inTyp] = true
		}
		if m.hasOut {
			needed[m.outTyp] = true
		}
	}
	var b strings.Builder
	b.WriteString("// Code generated by contractgen. DO NOT EDIT.\n\n")
	names := make([]string, 0, len(needed))
	for name := range needed {
		names = append(names, name)
	}
	// stable-ish order
	for _, name := range sortedStrings(names) {
		fields := structs[name]
		fmt.Fprintf(&b, "export type %s = {\n", name)
		for _, f := range fields {
			fmt.Fprintf(&b, "  %s: %s\n", f.name, f.typ)
		}
		b.WriteString("}\n\n")
	}
	return []byte(b.String())
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
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
		args := "signal?: AbortSignal"
		callBody := ""
		if m.hasIn {
			args = "input: contract." + m.inTyp + ", signal?: AbortSignal"
			callBody = ", body: input"
		}
		ret := "void"
		if m.hasOut {
			ret = "contract." + m.outTyp
		}
		fmt.Fprintf(&b, "export function %s(%s) {\n", fn, args)
		fmt.Fprintf(&b, "  return call<%s>({ method: %q, path: %q%s }, signal)\n}\n\n", ret, m.route.method, "/api"+m.route.path, callBody)
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

func modulePath(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			mod := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			if mod == "" {
				return "", fmt.Errorf("empty module path in go.mod")
			}
			return mod, nil
		}
	}
	return "", fmt.Errorf("module directive not found in go.mod")
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
