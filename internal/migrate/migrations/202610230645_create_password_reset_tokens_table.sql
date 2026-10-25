-- +goose Up
CREATE TABLE password_reset_tokens (
    id         bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    user_id    bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at    timestamptz
);

CREATE UNIQUE INDEX password_reset_tokens_token_hash_unique ON password_reset_tokens (token_hash);
CREATE INDEX password_reset_tokens_user_id_idx ON password_reset_tokens (user_id);

CREATE TRIGGER password_reset_tokens_set_updated_at
    BEFORE UPDATE ON password_reset_tokens
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE password_reset_tokens IS '密码重置令牌';
COMMENT ON COLUMN password_reset_tokens.id IS '主键';
COMMENT ON COLUMN password_reset_tokens.created_at IS '创建时间';
COMMENT ON COLUMN password_reset_tokens.updated_at IS '更新时间';
COMMENT ON COLUMN password_reset_tokens.user_id IS '用户 id';
COMMENT ON COLUMN password_reset_tokens.token_hash IS '重置令牌哈希';
COMMENT ON COLUMN password_reset_tokens.expires_at IS '过期时间';
COMMENT ON COLUMN password_reset_tokens.used_at IS '使用时间；空表示未使用';

-- +goose Down
DROP TABLE IF EXISTS password_reset_tokens;
