-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION set_updated_at() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END
$$;
-- +goose StatementEnd

COMMENT ON FUNCTION set_updated_at() IS '更新行时把 updated_at 写为事务时刻';

-- +goose Down
DROP FUNCTION IF EXISTS set_updated_at();
