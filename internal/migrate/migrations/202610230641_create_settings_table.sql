-- +goose Up
CREATE TABLE settings (
    name       text PRIMARY KEY,
    payload    jsonb NOT NULL DEFAULT '{}'::jsonb,
    blob       bytea,
    version    bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER settings_set_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE settings IS '系统功能设置，按 name 一行';
COMMENT ON COLUMN settings.name IS '设置组名';
COMMENT ON COLUMN settings.payload IS '该组 JSON 内容';
COMMENT ON COLUMN settings.blob IS '可选二进制附件';
COMMENT ON COLUMN settings.version IS '乐观版本，每次更新加一';
COMMENT ON COLUMN settings.created_at IS '创建时间';
COMMENT ON COLUMN settings.updated_at IS '更新时间';

-- +goose Down
DROP TABLE IF EXISTS settings;
