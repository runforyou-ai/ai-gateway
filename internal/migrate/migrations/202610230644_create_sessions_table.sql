-- +goose Up
CREATE TABLE sessions (
    id         bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    user_id    bigint NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX sessions_token_hash_unique ON sessions (token_hash);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TRIGGER sessions_set_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE sessions IS '管理台 Cookie 会话';
COMMENT ON COLUMN sessions.id IS '主键';
COMMENT ON COLUMN sessions.created_at IS '创建时间';
COMMENT ON COLUMN sessions.updated_at IS '更新时间';
COMMENT ON COLUMN sessions.user_id IS '用户 id';
COMMENT ON COLUMN sessions.token_hash IS 'Cookie 会话令牌哈希';
COMMENT ON COLUMN sessions.expires_at IS '过期时间（可滑动续期）';

-- +goose Down
DROP TABLE IF EXISTS sessions;
