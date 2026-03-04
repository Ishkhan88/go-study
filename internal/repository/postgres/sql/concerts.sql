-- name: CreateConcert :one
INSERT INTO concerts (title, date, tickets_total, tickets_left, genre, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, title, date, tickets_total, tickets_left, genre, description, created_at;

-- name: GetConcertByID :one
SELECT id, title, date, tickets_total, tickets_left, genre, description, created_at
FROM concerts
WHERE id = $1;

-- name: ListConcerts :many
SELECT id, title, date, tickets_total, tickets_left, genre, description, created_at
FROM concerts
ORDER BY date
LIMIT $1 OFFSET $2;

-- name: DecreaseTicketsLeft :exec
UPDATE concerts
SET tickets_left = tickets_left - $2
WHERE id = $1 AND tickets_left >= $2;