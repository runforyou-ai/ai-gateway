# 开发约定

## 硬约束

- 公共纯函数优先 `github.com/runforyou-ai/support`。
- 业务与运行日志只用 `log/slog`；控制台按行纯文本。
- 构建、生成、本地联调按 `Makefile`（`dev backend` / `dev frontend` / `build`）；`dev backend` 与 `build` 会跑契约生成；停后端用 `go run . stop`。
- 每表必须有 `created_at`、`updated_at`；列顺序紧跟主键（`id` 或 `name`）之后，并挂 `set_updated_at()` 触发器；model 字段顺序与迁移一致。

## 分层

```
cmd/                         # 无业务
internal/contract/           # 唯一契约源（//contract:route）
internal/tools/contractgen/  # 产物禁止手改
internal/direct/             # ops：调 actions、映射错误
internal/api/                # gen 的 /api；手写会话与 /v1
internal/actions/<域>/       # 不得 import contract/api/direct
internal/model/  internal/domain/  internal/provider/
internal/sysdoc/             # settings/state 按名读写
internal/integrationtest/
```

- `model` / `domain` / `provider` 不依赖上层。
- 契约 `auth` 仅 `public` | `user` | `admin`。
- `/v1/*` 不走契约生成。
- 管理台 Cookie 会话（`login` / `logout` / `me`）因 Set-Cookie 与契约 JSON 形状不合，手写在 `internal/api/session.go`，不进 `contract.Backend`。

## 配置边界

- yaml：基础设施（监听、Postgres、日志目录、邮件 Resend 等）；`init` 可读 `.env`，`serve` / `reload` 只读 yaml。
- 管理台可改项进 `settings`，运行单例进 `state`；二者经 `sysdoc` 按名读写，组名与 payload 结构集中定义。
- 可管理台改的项不进 yaml；连接串、邮件密钥等不进 `settings`。
- 异步队列（未来 jobq，驱动可选）与实时（jetcast）按需再接入，不作为启动硬依赖。

## 认证

- 管理台 Cookie 会话；`/v1` 用 Virtual Key，二者分离。
- 联调管理员：用户名 `Admin`，邮箱 `admin@example.com`，密码 `12345678`。

## `/api` 响应

- 成功：直接业务 JSON 或 204，不包 `{ data }`。
- 失败：`{ "message": "...", "fields": { "email": "required" } }`（`fields` 可选）；无单独业务 `code`。
- 500 的 `message` 对前端固定通用文案，细节只打 slog。
- `/v1` 错误保持协议原生形状。

## 其他

- 新迁移：`YYYYMMDDHHMM_描述.sql`。
- 前端优先 Refine / refine-ui / shadcn 官方写法；表单参考 `frontend/src/pages/profile.tsx`；实时若接入则走 jetcast，不做 realtime 契约生成。
- 集成测试：`internal/integrationtest/`；`TEST_POSTGRES_*`（至少 `HOST`），库名以 `_test` 结尾，未设则 Skip。
