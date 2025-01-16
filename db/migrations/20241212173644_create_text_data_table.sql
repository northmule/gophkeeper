-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.text_data (
      id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
      value text NOT NULL,
      "uuid" uuid NOT NULL,
      "name" varchar(300) NOT NULL,
      CONSTRAINT text_data_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS text_data;
-- +goose StatementEnd
