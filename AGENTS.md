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
| 认证 | JWT；`users` + `sessions`；`role` 枚举仅 `admin` / `user` |
| 日志 | 标准库 `log/slog`；控制台按行输出；文件轮转写入日志目录 |
| 辅助 | `github.com/runforyou-ai/support` |

打包与生成用 `make`（含 `make generate`）。

## 配置与进程

- 默认配置：`$XDG_CONFIG_HOME/lite-ai-gateway/config.yaml`（未设 XDG 时为 `~/.config/...`）。
- `init` 写入该路径（可参数指定路径）；`serve` 只从该默认路径读，或 `-f` 指定；找不到则报错退出。
- 开发：允许环境变量覆盖 yaml 缺省；生产：忽略环境变量，只认 yaml。
- NATS：yaml `nats.url` / 开发可用 `NATS_URL`；空=嵌入，非空=独立或集群。
- pid 默认：与配置同目录的 `lite-ai-gateway.pid`。

### 启动配置 vs 系统 settings / state

参考 luway-cloud：启动只读 yaml；运行期系统级配置与单例状态进 Postgres。

| 层 | 存哪 | 内容 | 谁改 |
|---|---|---|---|
| 启动配置 | yaml（及开发期 env） | 监听地址、Postgres、NATS、日志目录、JWT 密钥等基础设施 | 运维 / `init`；进程 `reload` 可重载 yaml |
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
internal/auth/                # JWT、会话原语
internal/config/  internal/nats/
```

依赖规则：

- `actions` 不得 import `contract`、`api`、`direct`。
- `model`、`domain`、`provider` 不依赖上层。
- 契约方法 `auth` 仅：`public` | `user` | `admin`（无工作区、无自定义 perm）。
- `/v1/*` 不走契约生成（流式与多协议），在 `api` 手写。
- 不强制全局 `repo` 层；查询重再在对应 action 包内拆文件。
- 公共纯函数优先 support 或小包，避免「utils 垃圾场」。

### 契约与生成

- 在 `contract` 声明方法及 `//contract:route ...`；改完后执行 `make generate`。
- 生成物：`direct` 分发、`api` 路由、`frontend/src/api/generated/`（`contract.ts`、`operations.ts`）。
- 生成文件标 `DO NOT EDIT`，禁止手改；只改契约或生成器后重新生成。
- 手写 `direct` ops：调用 action、做类型/错误映射；业务校验（存在性、唯一性等）在 action。

## 数据库

- 迁移文件名：`YYYYMMDDHHMM_描述.sql`（例：`202610230611_create_users_table.sql`）。
- 普通业务表：`id`（标准自增）、`created_at`、`updated_at`；用统一的 `set_updated_at()` 触发器。
- `settings` / `state`：见上文文档组模型（`name` 主键，无自增 `id`）。
- 一表一 model；字段与 bun tag 与迁移一致。
- 首个注册用户：`role=admin`；之后默认 `user`。无自定义角色、不做复杂 RBAC。

## 前端

- 管理台经 `@/api`（含 generated operations）调 `/api`；可用薄 Refine `dataProvider` 封装。
- 实时走 jetcast（浏览器直连 NATS WebSocket），不做 realtime 契约代码生成。
- 样式用既有 shadcn + Tailwind；不引入第二套 UI 体系。
- 构建产物嵌入 Go；本地 `refine dev` 可独立联调。

## 文档与杂项

- 产品行为与架构以 `README.md` 为准；本文件只约束实现方式。
- 不提交密钥；`.env` 仅本地。`./data/` 下运行时数据不入库（见 `.gitignore`）。
