-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscribers(
	id int NOT NULL,
	is_valid boolean NOT NULL,
	expires_at timestamptz NOT NULL
);

-- +goose StatementEnd

INSERT INTO subscribers (id, expires_at, is_valid) VALUES
(1, '2025-05-10 15:00:00', true),
(2, '2025-05-01 12:00:00', true),
(3, '2025-04-30 06:00:00', true);


-- +goose Down
-- +goose StatementBegin
DROP TABLE subscribers;
-- +goose StatementEnd