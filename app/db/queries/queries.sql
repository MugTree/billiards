-- name: SelectPeople :many
SELECT * FROM people;

-- name: InsertPerson :exec
INSERT INTO people (name) VALUES (?);
