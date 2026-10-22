-- +goose Up
CREATE FUNCTION set_updated_at() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$;

-- +goose Down
DROP FUNCTION IF EXISTS set_updated_at();
