-- name: CreateBooking :one
INSERT INTO bookings (user_id, concert_id, status, quantity, note)
VALUES ($1, $2, 'pending', $3, $4)
RETURNING id, user_id, concert_id, status, quantity, note, created_at;

-- name: GetBookingByID :one
SELECT id, user_id, concert_id, status, quantity, note, created_at
FROM bookings
WHERE id = $1;

-- name: SetBookingStatus :exec
UPDATE bookings
SET status = $2
WHERE id = $1;

-- name: ListBookingsByUser :many
SELECT id, user_id, concert_id, status, quantity, note, created_at
FROM bookings
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;