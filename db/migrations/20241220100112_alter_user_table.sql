-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ADD private_client_key text NULL;
ALTER TABLE public.users ADD private_server_key text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
