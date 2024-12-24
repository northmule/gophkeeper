-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ALTER COLUMN public_key SET DEFAULT 'n'::text;
ALTER TABLE public.users DROP COLUMN private_server_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
