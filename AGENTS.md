# 开发约定

面向本仓库的 AI Agent / 协作者。只记录必须遵守的约定，细节以代码与 README 为准。

## 范围

- 只做用户明确要求的改动；发现隐藏 bug 只提出建议，不擅自改。
- 仅构建与部署 **Linux**；不为其他平台加构建目标。
- 用户可见说明与提交说明用中文亦可；代码标识符用英文。

## 技术栈

| 层 | 选型 |
|---|---|
| CLI / 进程 | Cobra：`init`、`serve`、`stop`、`restart`、`reload` |
| HTTP | 标准库 `net/http`；`/v1/*` AI 协议（手写），`/api/*` 管理 API（契约生成） |
| 契约生成 | `internal/tools/contractgen`：由 `contract.Backend` 生成 direct 分发、HTTP 路由、前端 TS |
| 协议 | 仅平台预设供应商；透传/转换矩阵写死在预设中；跨协议经 RelayKit；不支持自定义供应商 |
| DB | Postgres；goose 迁移；bun 访问；每表对应一个 model |
| 消息 | NATS JetStream；jetcast 广播；jetq 任务；`nats.url` 空则嵌入，非空则外连 |
| 前端 | Refine + shadcn + Tailwind；产物 `embed` 进二进制 |
| 认证 | `/api`：HttpOnly Cookie + `sessions`；`/v1`：API Key（后续）；`users.role` 仅 `admin` / `user` |
| 日志 | 标准库 `log/slog`；控制台按行输出；文件轮转写入日志目录 |
| 辅助 | `github.com/runforyou-ai/support` |

Make：`make dev`（前后端）；`make build`（当前系统二进制）；`make build amd64|arm64`（Linux 交叉编译）。
Docker（生产单机）：先 `make build amd64|arm64`（与宿主机 Docker 架构一致），再 `docker compose up --build`；仅 Postgres + app，NATS 内嵌；镜像 COPY 对应 Linux 二进制；挂载仓库根目录 `config.yaml`、日志与 JetStream 数据目录。
CI：GitHub Actions（PR / `main`）跑契约生成校验、`go vet`/`test`、前端构建，以及 compose 冒烟（`/health`）；`main` 用 buildx 推送 `linux/amd64`+`linux/arm64` 多架构镜像到 `ghcr.io/runforyou-ai/lite-ai-gateway`（`latest` 与 commit SHA）。
CD：推送 tag `v*` 触发 release 工作流，发布 GitHub Release，并推送镜像 tag（`vX.Y.Z`、`X.Y.Z`、`latest`）。

## 配置与进程

- `make dev` / `serve` 读默认 XDG 配置：`$XDG_CONFIG_HOME/lite-ai-gateway/config.yaml`（未设 XDG 时为 `~/.config/...`）；需自行 `go run . init`。
- `init`：内置默认值；若当前目录有 `.env`，用其中变量覆盖后再写 yaml（`-o` 可指定路径）。
- `serve` / `reload`：只读 yaml（默认路径或 `-f`）；不读环境变量、不读 `.env`；找不到 yaml 则报错退出。

- NATS：只认 yaml `nats.url`；空=嵌入，非空=独立或集群。
- pid 默认：与配置同目录的 `lite-ai-gateway.pid`。
- yaml 里的 `env`（`development` / `production`）只影响日志等默认路径，不改变「是否读环境变量」。

### 启动配置 vs 系统 settings / state

参考 luway-cloud：启动只读 yaml；运行期系统级配置与单例状态进 Postgres。

| 层 | 存哪 | 内容 | 谁改 |
|---|---|---|---|
| 启动配置 | yaml（`init` 时可从 `.env` 生成） | 监听地址、Postgres、NATS、日志目录等基础设施 | 运维 / `init`；进程 `reload` 可重载 yaml |
| `settings` | 表，按 `name` 一组 JSON（可附 `blob`） | 管理员可改的系统设置（如站点名、是否开放注册等） | 管理台 / admin API |
| `state` | 表，按 `name` 一组 JSON | 部署/运行单例状态（如已安装标记、运行态凭据等），非日常表单 | 安装与内部逻辑写入 |

约定：

- 不要把可在管理台改的项放进 yaml；不要把连接串等基础设施放进 `settings`。
- `settings` / `state`：`name` 主键 + `payload` jsonb + `version`（乐观锁，更新 +1）+ `created_at` / `updated_at`；`settings` 可有可选 `blob`。此二表**不以自增 `id` 为主键**（文档组模型例外）。
- 读写封装在 `internal/sysdoc`（按名 Get/Put/Lock）；各组 payload 结构与组名常量放一处，禁止散落拼 JSON。
- 热路径可读进程内快照（类似 luway `DeploymentState`）：启动与 settings/state 变更后 `Reload`；多实例下他机可下次刷新或经广播失效。

## 日志

- 使用 `log/slog`；禁止用 `fmt`/`log` 打业务与运行日志。
- 控制台：按行输出纯文本（无表情、无 ANSI 花活）。
- 同时轮转写入文件：开发默认 `./data/logs/`；生产默认 `/var/log/lite-ai-gateway/`；yaml 中配置日志目录时可覆盖默认路径。
- 文件轮转（按大小/保留份数等）在实现里统一封装，业务代码只调 slog。

## 后端分层（Action + 契约）

参考 luway-cloud，命名更直白；单向依赖：

```
cmd/                          # Cobra，不写业务
internal/contract/            # Backend 接口 + DTO + //contract:route（唯一契约源）
internal/tools/contractgen/   # 生成器；禁止手改其产物
internal/direct/              # gen 认证分发 + 手写 ops（调 actions、映射 DTO/错误）
internal/api/                 # gen 的 /api 路由；手写 /v1 透传与转换
internal/actions/<域>/        # 业务用例；可用 bun，不依赖 contract/api/direct
internal/model/               # 与表一一对应的 bun 结构体
internal/domain/              # 枚举与领域常量（role、协议等）
internal/provider/            # 预设供应商及写死的协议透传/转换矩阵
internal/sysdoc/              # settings/state 按名文档读写
internal/auth/                # Cookie 会话、Principal 原语
internal/config/  internal/nats/
```

依赖规则：

- `actions` 不得 import `contract`、`api`、`direct`。
- `model`、`domain`、`provider` 不依赖上层。
- 契约方法 `auth` 仅：`public` | `user` | `admin`（无工作区、无自定义 perm）。
- `/v1/*` 不走契约生成（流式与多协议），在 `api` 手写。
- 不强制全局 `repo` 层；查询重再在对应 action 包内拆文件。
- 公共纯函数优先 support 或小包，避免「utils 垃圾场」。

### 认证与注册

- **无单独安装向导**。通用 `POST /api/register`：系统中第一个注册用户自动为 `admin`，之后注册为 `user`（需 `settings.site.registrationOpen`；无人时始终允许注册）。
- **本地/联调默认管理员**（Agent 或脚本自动创建账号时必须用这对凭据，勿另起邮箱、用户名或密码）：
  - 用户名：`Admin`
  - 邮箱：`admin@example.com`
  - 密码：`12345678`
- `users.name` 为展示用用户名（注册必填）；登录标识仍是 `email`。
- **管理台 `/api`**：浏览器 Cookie 会话（`HttpOnly`；生产加 `Secure`；`SameSite=Lax` 或按部署调整）。Cookie 存不透明会话令牌，服务端 `sessions` 表为权威（可吊销、滑动过期、多设备）。
- **不使用 JWT** 做管理台登录态。
- 登录：校验账号后插入 `sessions`（存令牌哈希、`user_id`、`expires_at` 等），Set-Cookie；登出：删会话行并清 Cookie。
- 鉴权：读 Cookie → 查 `sessions`（有效且未过期）→ 加载 `users`；`user`/`admin` 按 `role` 判断。活跃请求可滑动延长 `expires_at`。
- **对外 `/v1`**：使用 API Key（或其它机器凭证），不用浏览器 Cookie；与管理台会话分离。
- 同域嵌入前端优先，避免跨域 Cookie/CSRF 复杂度；若将来跨域再单独补 CORS/CSRF。
- `serve` / `reload` 只读 yaml；仅 `init` 可读当前目录 `.env`。

### 契约与生成

- 在 `contract` 声明方法及 `//contract:route ...`；改完后执行 `make generate`。
- 生成物：`direct` 分发、`api` 路由、`frontend/src/api/generated/`（`contract.ts`、`operations.ts`）。
- 生成文件标 `DO NOT EDIT`，禁止手改；只改契约或生成器后重新生成。
- 手写 `direct` ops：调用 action、做类型/错误映射；业务校验（存在性、唯一性等）在 action。

### `/api` 错误与成功约定

成功：直接返回业务 JSON（或 204 无体），**不**包一层 `{ data: ... }`。

失败：用 **HTTP 状态码** 区分错误大类（前端按 `response.status` 分支）；body 统一为：

```json
{
  "message": "可读说明",
  "fields": { "email": "required" }
}
```

- `message`：给人看的摘要；第一期可用中文或英文固定文案，不做 i18n 键体系。
- `fields`：可选；仅 422 等需要表单高亮时使用；字段名 → 短错误码（如 `required`、`invalid`、`taken`）。
- 不返回单独的业务 `code` 字段（与状态码一一对应时属冗余）；若将来同一状态码需多种语义再补。
- 未预期错误：`500`；`message` 对前端固定为通用文案，细节只打 slog。

| HTTP | 含义 |
|---|---|
| 400 | 无法解析的请求（坏 JSON、缺 body 等） |
| 422 | 请求已解析但业务校验失败（可带 `fields`） |
| 401 | 未登录或令牌无效 |
| 403 | 已登录但无权限 |
| 404 | 资源不存在 |
| 409 | 冲突（唯一约束、状态不允许） |
| 429 | 限流（可选后续加 `retryAfter`） |
| 500 | 未分类服务端错误 |

实现约定：

- Action 返回可识别的领域错误；`direct` 映射为上述状态码与 `message`/`fields`；`api` 层统一编码。
- 前端 `call()`：非 2xx 按状态码分支，解析 `message`/`fields`；表单吃 `fields`，其余 Toast 展示 `message`。
- `/v1/*` AI 协议错误保持上游/协议原生形状，**不**套用本表（管理台 `/api` 专用）。

## 数据库

- 迁移文件名：`YYYYMMDDHHMM_描述.sql`（例：`202610230611_create_users_table.sql`）。
- 普通业务表：`id`（标准自增）、`created_at`、`updated_at`；用统一的 `set_updated_at()` 触发器。
- `settings` / `state`：见上文文档组模型（`name` 主键，无自增 `id`）。
- 一表一 model；字段与 bun tag 与迁移一致。
- 首个注册用户：`role=admin`；之后默认 `user`。无自定义角色、不做复杂 RBAC。
- `users`：`name`（展示名）、`email`（登录）、`password`、`role`。

## 前端

- 管理台经 `@/api`（含 generated operations）调 `/api`；可用薄 Refine `dataProvider` 封装。
- 实时走 jetcast（浏览器直连 NATS WebSocket），不做 realtime 契约代码生成。
- 样式用既有 shadcn + Tailwind；不引入第二套 UI 体系。
- 构建产物嵌入 Go；本地 `refine dev` 可独立联调。
- **优先用官方能力**：页面、布局、表单、表格、按钮等尽可能用 Refine / refine-ui / shadcn 自带组件与模板写法（含其默认间距与样式）；不要手写花哨布局或大量自定义 class 堆样式。确需定制时再局部覆盖，并保持简洁。表单参考 `frontend/src/pages/profile.tsx`。

## 文档与杂项

- 产品行为与架构以 `README.md` 为准；本文件只约束实现方式。
- 不提交密钥；`.env` 仅本地。`./data/` 下运行时数据不入库（见 `.gitignore`）。
