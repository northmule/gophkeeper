-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.file_data (
          id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
          "name" varchar(300) NOT NULL,
          "uuid" uuid NOT NULL,
          mime_type varchar(50) NOT NULL,
          "path" varchar(300) NOT NULL,
          "extension" varchar(10) NOT NULL,
          file_name varchar(300) NOT NULL,
          "storage" varchar(100) NOT NULL,
          uploaded bool NOT NULL,
          path_tmp varchar(300) NULL,
          "size" int8 NOT NULL,
          CONSTRAINT binary_data_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS file_data;
-- +goose StatementEnd
