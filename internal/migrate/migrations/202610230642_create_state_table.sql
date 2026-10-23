-- +goose Up
CREATE TABLE state (
    name       text PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    payload    jsonb NOT NULL DEFAULT '{}'::jsonb,
    version    bigint NOT NULL DEFAULT 1
);

CREATE TRIGGER state_set_updated_at
    BEFORE UPDATE ON state
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE state IS '系统单例状态，按 name 一行';
COMMENT ON COLUMN state.name IS '状态组名';
COMMENT ON COLUMN state.created_at IS '创建时间';
COMMENT ON COLUMN state.updated_at IS '更新时间';
COMMENT ON COLUMN state.payload IS '该组 JSON 内容';
COMMENT ON COLUMN state.version IS '乐观版本，每次更新加一';

-- +goose Down
DROP TABLE IF EXISTS state;
