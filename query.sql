-- name: CreateVisit :one
INSERT INTO visit (start_time_unix, length_second)
VALUES (?,?)
RETURNING id
;

-- name: ListVisit :many
SELECT * FROM visit ORDER BY start_time_unix
;
