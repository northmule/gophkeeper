-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ADD public_key text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
