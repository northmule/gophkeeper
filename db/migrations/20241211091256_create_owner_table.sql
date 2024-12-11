-- +goose Up
-- +goose StatementBegin
CREATE TABLE public."owner" (
        id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
        user_uuid uuid NOT NULL,
        data_type varchar(50) NOT NULL,
        data_uuid uuid NOT NULL,
        CONSTRAINT owner_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "owner";
-- +goose StatementEnd
