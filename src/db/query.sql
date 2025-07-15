-- name: CreateUser :execresult
INSERT INTO users (
    name, email, google_id, avatar_url, role
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?
  AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
  AND deleted_at IS NULL;

-- name: GetUserByGoogleID :one
SELECT * FROM users
WHERE google_id = ?
  AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY id;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- name: CreateReservation :execresult
INSERT INTO reservations (
    user_id, title, start_time, end_time, status
) VALUES (
    ?, ?, ?, ?, 'active'
);

-- name: GetReservationByID :one
SELECT * FROM reservations
WHERE id = ?;

-- name: ListReservationsByUserID :many
SELECT * FROM reservations
WHERE user_id = ?
ORDER BY start_time;

-- name: ListReservationsByMonth :many
SELECT r.*, u.name
FROM reservations r
JOIN users u ON r.user_id = u.id
WHERE DATE_FORMAT(r.start_time, '%Y-%m') = DATE_FORMAT(?, '%Y-%m');

-- name: ListReservationsByWeek :many
SELECT r.*, u.name
FROM reservations r
JOIN users u ON r.user_id = u.id
WHERE r.start_time >= ? AND r.end_time <= ?;

-- name: ListReservationsByDate :many
SELECT r.*, u.name
FROM reservations r
JOIN users u ON r.user_id = u.id
WHERE DATE(r.start_time) = DATE(?);

-- name: UpdateReservationByID :exec
UPDATE reservations
SET title = ?, start_time = ?, end_time = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteReservationByID :exec
DELETE FROM reservations
WHERE id = ?;
