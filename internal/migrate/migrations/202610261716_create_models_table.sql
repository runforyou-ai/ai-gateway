-- +goose Up
CREATE TABLE models (
    id              bigserial PRIMARY KEY,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    name            text NOT NULL,
    display_name    text NOT NULL,
    type            text NOT NULL,
    disabled        boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX models_name_unique ON models (name);

CREATE TRIGGER models_set_updated_at
    BEFORE UPDATE ON models
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE models IS '对外模型目录（公开契约名）';
COMMENT ON COLUMN models.id IS '主键';
COMMENT ON COLUMN models.created_at IS '创建时间';
COMMENT ON COLUMN models.updated_at IS '更新时间';
COMMENT ON COLUMN models.name IS '对外模型 id（如 deepseek-flash）';
COMMENT ON COLUMN models.display_name IS '展示名称';
COMMENT ON COLUMN models.type IS '模型类型：chat 等；对外协议由代码按类型推导';
COMMENT ON COLUMN models.disabled IS '是否禁用；禁用后不可被路由选中';

-- +goose Down
DROP TABLE IF EXISTS models;
