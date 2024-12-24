-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ALTER COLUMN public_key SET DEFAULT 'n'::text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
