-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.meta_data
(
    id         int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
    meta_name  varchar(50)                       NOT NULL,
    meta_value jsonb                             NOT NULL,
    data_uuid  uuid                              NOT NULL,
    CONSTRAINT meta_data_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meta_data;
-- +goose StatementEnd
