-- name: CreateMember :one
INSERT INTO members (
    first_name,
    last_name,
    birth_date,
    email
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    first_name,
    last_name,
    birth_date,
    email,
    created_at;