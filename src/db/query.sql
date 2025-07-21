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
    ?, ?, ?, ?, 'confirmed'
);

-- name: GetReservationLastInserted :one
SELECT * FROM reservations
WHERE id = LAST_INSERT_ID();


-- name: GetReservationByID :one
SELECT * FROM reservations
WHERE id = ?;

-- name: ListReservationsByUserID :many
SELECT * FROM reservations
WHERE status = 'confirmed'
  AND user_id = ?
ORDER BY start_time;

-- name: ListReservationsByMonth :many
SELECT r.*, u.name as user_name
FROM reservations AS r
JOIN users AS u ON r.user_id = u.id AND u.deleted_at IS NULL
WHERE
  r.status = 'confirmed'
  AND r.start_time < ?  -- 翌月の初日
  AND r.end_time >= ? -- 月の初日
ORDER BY
  r.start_time;

-- name: ListReservationsByWeek :many
SELECT r.*, u.name as user_name
FROM reservations AS r
JOIN users AS u ON r.user_id = u.id AND u.deleted_at IS NULL
WHERE
  r.status = 'confirmed'
  AND r.start_time < sqlc.arg(EndTime)  
  AND r.end_time >= sqlc.arg(StartTime) 
ORDER BY
  r.start_time;

-- name: ListReservationsByDate :many
SELECT r.*, u.name as user_name
FROM reservations AS r
JOIN users AS u ON r.user_id = u.id AND u.deleted_at IS NULL
WHERE
  r.status = 'confirmed'
  AND r.start_time < ? 
  AND r.end_time >= ? 
ORDER BY
  r.start_time;

-- name: UpdateReservationByID :exec
UPDATE reservations
SET title = ?, start_time = ?, end_time = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteReservationByID :exec
DELETE FROM reservations
WHERE user_id = ? 
  AND id = ?;

-- name: CanceledReservationByID :exec
UPDATE reservations
SET
  status = 'canceled',
  updated_at = CURRENT_TIMESTAMP
WHERE
  user_id = ? AND id = ?;

-- name: CheckOverlappingReservation :one
SELECT COUNT(*) FROM reservations
WHERE status = 'confirmed'
  AND start_time < ?
  AND end_time > ?


