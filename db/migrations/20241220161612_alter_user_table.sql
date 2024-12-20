-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE public.users ALTER COLUMN private_client_key SET DEFAULT 'n';
ALTER TABLE public.users ALTER COLUMN private_server_key SET DEFAULT 'n';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
