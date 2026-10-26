-- +goose Up
CREATE TABLE virtual_keys (
    id            bigserial PRIMARY KEY,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    user_id       bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name          text NOT NULL,
    prefix        text NOT NULL,
    token_hash    text NOT NULL,
    disabled      boolean NOT NULL DEFAULT false,
    max_budget    numeric(20, 8),
    spent         numeric(20, 8) NOT NULL DEFAULT 0,
    budget_reset  text NOT NULL DEFAULT '',
    expires_in    text NOT NULL DEFAULT '',
    expires_at    timestamptz,
    last_used_at  timestamptz
);

CREATE UNIQUE INDEX virtual_keys_token_hash_unique ON virtual_keys (token_hash);
CREATE INDEX virtual_keys_user_id_idx ON virtual_keys (user_id);

CREATE TRIGGER virtual_keys_set_updated_at
    BEFORE UPDATE ON virtual_keys
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE virtual_keys IS '调用 /v1 的虚拟密钥';
COMMENT ON COLUMN virtual_keys.id IS '主键';
COMMENT ON COLUMN virtual_keys.created_at IS '创建时间';
COMMENT ON COLUMN virtual_keys.updated_at IS '更新时间';
COMMENT ON COLUMN virtual_keys.user_id IS '所属用户 id';
COMMENT ON COLUMN virtual_keys.name IS '显示名称';
COMMENT ON COLUMN virtual_keys.prefix IS '明文前缀（便于识别）';
COMMENT ON COLUMN virtual_keys.token_hash IS '完整令牌哈希';
COMMENT ON COLUMN virtual_keys.disabled IS '是否禁用';
COMMENT ON COLUMN virtual_keys.max_budget IS '消费上限；空表示不限';
COMMENT ON COLUMN virtual_keys.spent IS '已消费金额';
COMMENT ON COLUMN virtual_keys.budget_reset IS '额度重置周期：空/daily/weekly/monthly';
COMMENT ON COLUMN virtual_keys.expires_in IS '过期时长编码：空/1h/1d/3d/7d/30d/90d/180d/365d';
COMMENT ON COLUMN virtual_keys.expires_at IS '过期时间；空表示不过期';
COMMENT ON COLUMN virtual_keys.last_used_at IS '最近使用时间';

-- +goose Down
DROP TABLE IF EXISTS virtual_keys;
