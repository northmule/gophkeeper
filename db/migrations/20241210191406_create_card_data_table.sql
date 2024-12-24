-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.card_data (
      id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
      value jsonb NOT NULL,
      object_type varchar(50) NOT NULL,
      "name" varchar(300) NOT NULL,
      CONSTRAINT card_data_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS card_data;
-- +goose StatementEnd
