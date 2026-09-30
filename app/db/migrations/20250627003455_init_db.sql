-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS people
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL DEFAULT ''
);
-- +goose StatementEnd

