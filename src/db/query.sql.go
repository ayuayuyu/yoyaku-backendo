package db

import (
	"context"
	"database/sql"
	"time"
)

const canceledReservationByID = `-- name: CanceledReservationByID :exec
UPDATE reservations
SET
  status = 'canceled',
  updated_at = CURRENT_TIMESTAMP
WHERE
  user_id = ? AND id = ?
`

type CanceledReservationByIDParams struct {
	UserID uint64 `json:"user_id"`
	ID     uint64 `json:"id"`
}

func (q *Queries) CanceledReservationByID(ctx context.Context, arg CanceledReservationByIDParams) error {
	_, err := q.db.ExecContext(ctx, canceledReservationByID, arg.UserID, arg.ID)
	return err
}

const checkOverlappingReservation = `-- name: CheckOverlappingReservation :one
SELECT COUNT(*) FROM reservations
WHERE status = 'confirmed'
  AND start_time < ?
  AND end_time > ?
`

type CheckOverlappingReservationParams struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (q *Queries) CheckOverlappingReservation(ctx context.Context, arg CheckOverlappingReservationParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkOverlappingReservation, arg.StartTime, arg.EndTime)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createReservation = `-- name: CreateReservation :execresult
INSERT INTO reservations (
    user_id, title, start_time, end_time, status
) VALUES (
    ?, ?, ?, ?, 'confirmed'
)
`

type CreateReservationParams struct {
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (q *Queries) CreateReservation(ctx context.Context, arg CreateReservationParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createReservation,
		arg.UserID,
		arg.Title,
		arg.StartTime,
		arg.EndTime,
	)
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (
    name, email, google_id, avatar_url, role
) VALUES (
    ?, ?, ?, ?, ?
)
`

type CreateUserParams struct {
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	GoogleID  string         `json:"google_id"`
	AvatarUrl sql.NullString `json:"avatar_url"`
	Role      string         `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.GoogleID,
		arg.AvatarUrl,
		arg.Role,
	)
}

const deleteReservationByID = `-- name: DeleteReservationByID :exec
DELETE FROM reservations
WHERE user_id = ? 
  AND id = ?
`

type DeleteReservationByIDParams struct {
	UserID uint64 `json:"user_id"`
	ID     uint64 `json:"id"`
}

func (q *Queries) DeleteReservationByID(ctx context.Context, arg DeleteReservationByIDParams) error {
	_, err := q.db.ExecContext(ctx, deleteReservationByID, arg.UserID, arg.ID)
	return err
}

const getReservationByID = `-- name: GetReservationByID :one
SELECT id, user_id, title, start_time, end_time, status, created_at, updated_at FROM reservations
WHERE id = ?
`

func (q *Queries) GetReservationByID(ctx context.Context, id uint64) (Reservation, error) {
	row := q.db.QueryRowContext(ctx, getReservationByID, id)
	var i Reservation
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.StartTime,
		&i.EndTime,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getReservationLastInserted = `-- name: GetReservationLastInserted :one
SELECT id, user_id, title, start_time, end_time, status, created_at, updated_at FROM reservations
WHERE id = LAST_INSERT_ID()
`

func (q *Queries) GetReservationLastInserted(ctx context.Context) (Reservation, error) {
	row := q.db.QueryRowContext(ctx, getReservationLastInserted)
	var i Reservation
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.StartTime,
		&i.EndTime,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, google_id, avatar_url, role, created_at, updated_at, deleted_at FROM users
WHERE email = ?
  AND deleted_at IS NULL
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.GoogleID,
		&i.AvatarUrl,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByGoogleID = `-- name: GetUserByGoogleID :one
SELECT id, name, email, google_id, avatar_url, role, created_at, updated_at, deleted_at FROM users
WHERE google_id = ?
  AND deleted_at IS NULL
`

func (q *Queries) GetUserByGoogleID(ctx context.Context, googleID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByGoogleID, googleID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.GoogleID,
		&i.AvatarUrl,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, google_id, avatar_url, role, created_at, updated_at, deleted_at FROM users
WHERE id = ?
  AND deleted_at IS NULL
`

func (q *Queries) GetUserByID(ctx context.Context, id uint64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.GoogleID,
		&i.AvatarUrl,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listReservationsByDate = `-- name: ListReservationsByDate :many
SELECT r.id, r.user_id, r.title, r.start_time, r.end_time, r.status, r.created_at, r.updated_at, u.name as user_name
FROM reservations AS r
JOIN users AS u ON r.user_id = u.id AND u.deleted_at IS NULL
WHERE
  r.status = 'confirmed'
  AND r.start_time < ? 
  AND r.end_time >= ? 
ORDER BY
  r.start_time
`

type ListReservationsByDateParams struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type ListReservationsByDateRow struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
}

func (q *Queries) ListReservationsByDate(ctx context.Context, arg ListReservationsByDateParams) ([]ListReservationsByDateRow, error) {
	rows, err := q.db.QueryContext(ctx, listReservationsByDate, arg.StartTime, arg.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListReservationsByDateRow
	for rows.Next() {
		var i ListReservationsByDateRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listReservationsByMonth = `-- name: ListReservationsByMonth :many
SELECT r.id, r.user_id, r.title, r.start_time, r.end_time, r.status, r.created_at, r.updated_at, u.name as user_name
FROM reservations AS r
JOIN users AS u ON r.user_id = u.id AND u.deleted_at IS NULL
WHERE
  r.status = 'confirmed'
  AND r.start_time < ?  -- 翌月の初日
  AND r.end_time >= ? -- 月の初日
ORDER BY
  r.start_time
`

type ListReservationsByMonthParams struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type ListReservationsByMonthRow struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
}

func (q *Queries) ListReservationsByMonth(ctx context.Context, arg ListReservationsByMonthParams) ([]ListReservationsByMonthRow, error) {
	rows, err := q.db.QueryContext(ctx, listReservationsByMonth, arg.StartTime, arg.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListReservationsByMonthRow
	for rows.Next() {
		var i ListReservationsByMonthRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listReservationsByUserID = `-- name: ListReservationsByUserID :many
SELECT id, user_id, title, start_time, end_time, status, created_at, updated_at FROM reservations
WHERE status = 'confirmed'
  AND user_id = ?
ORDER BY start_time
`

func (q *Queries) ListReservationsByUserID(ctx context.Context, userID uint64) ([]Reservation, error) {
	rows, err := q.db.QueryContext(ctx, listReservationsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reservation
	for rows.Next() {
		var i Reservation
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listReservationsByWeek = `-- name: ListReservationsByWeek :many
SELECT r.id, r.user_id, r.title, r.start_time, r.end_time, r.status, r.created_at, r.updated_at, u.name as user_name
FROM reservations AS r
JOIN users AS u ON r.user_id = u.id AND u.deleted_at IS NULL
WHERE
  r.status = 'confirmed'
  AND r.start_time < ?  
  AND r.end_time >= ? 
ORDER BY
  r.start_time
`

type ListReservationsByWeekParams struct {
	Endtime   time.Time `json:"endtime"`
	Starttime time.Time `json:"starttime"`
}

type ListReservationsByWeekRow struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
}

func (q *Queries) ListReservationsByWeek(ctx context.Context, arg ListReservationsByWeekParams) ([]ListReservationsByWeekRow, error) {
	rows, err := q.db.QueryContext(ctx, listReservationsByWeek, arg.Endtime, arg.Starttime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListReservationsByWeekRow
	for rows.Next() {
		var i ListReservationsByWeekRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, email, google_id, avatar_url, role, created_at, updated_at, deleted_at FROM users
WHERE deleted_at IS NULL
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.GoogleID,
			&i.AvatarUrl,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteUser = `-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?
`

func (q *Queries) SoftDeleteUser(ctx context.Context, id uint64) error {
	_, err := q.db.ExecContext(ctx, softDeleteUser, id)
	return err
}

const updateReservationByID = `-- name: UpdateReservationByID :exec
UPDATE reservations
SET title = ?, start_time = ?, end_time = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

type UpdateReservationByIDParams struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	ID        uint64    `json:"id"`
}

func (q *Queries) UpdateReservationByID(ctx context.Context, arg UpdateReservationByIDParams) error {
	_, err := q.db.ExecContext(ctx, updateReservationByID,
		arg.Title,
		arg.StartTime,
		arg.EndTime,
		arg.ID,
	)
	return err
}
