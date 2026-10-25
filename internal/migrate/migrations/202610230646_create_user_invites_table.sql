-- +goose Up
CREATE TABLE user_invites (
    id          bigserial PRIMARY KEY,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    email       text NOT NULL,
    role        text NOT NULL,
    token_hash  text NOT NULL,
    expires_at  timestamptz NOT NULL,
    used_at     timestamptz,
    invited_by  bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX user_invites_token_hash_unique ON user_invites (token_hash);
CREATE INDEX user_invites_email_idx ON user_invites (lower(email));

CREATE TRIGGER user_invites_set_updated_at
    BEFORE UPDATE ON user_invites
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE user_invites IS '用户邀请令牌';
COMMENT ON COLUMN user_invites.id IS '主键';
COMMENT ON COLUMN user_invites.created_at IS '创建时间';
COMMENT ON COLUMN user_invites.updated_at IS '更新时间';
COMMENT ON COLUMN user_invites.email IS '受邀邮箱';
COMMENT ON COLUMN user_invites.role IS '接受后角色：admin 或 user';
COMMENT ON COLUMN user_invites.token_hash IS '邀请令牌哈希';
COMMENT ON COLUMN user_invites.expires_at IS '过期时间';
COMMENT ON COLUMN user_invites.used_at IS '使用时间；空表示未使用';
COMMENT ON COLUMN user_invites.invited_by IS '发起邀请的用户 id';

-- +goose Down
DROP TABLE IF EXISTS user_invites;
