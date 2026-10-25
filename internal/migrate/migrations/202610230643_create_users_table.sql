-- +goose Up
CREATE TABLE users (
    id            bigserial PRIMARY KEY,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    name          text NOT NULL,
    email         text NOT NULL,
    password      text NOT NULL,
    role          text NOT NULL,
    disabled      boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX users_email_unique ON users (lower(email));

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE users IS '用户账号';
COMMENT ON COLUMN users.id IS '主键';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.name IS '用户名（展示用）';
COMMENT ON COLUMN users.email IS '登录邮箱';
COMMENT ON COLUMN users.password IS '密码哈希';
COMMENT ON COLUMN users.role IS '角色：admin 或 user';
COMMENT ON COLUMN users.disabled IS '是否禁用；禁用后不可登录';

-- +goose Down
DROP TABLE IF EXISTS users;
