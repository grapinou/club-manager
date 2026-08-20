-- name: CreatePerson :one
INSERT INTO persons (
    first_name,
    last_name,
    birth_date,
    phone_number,
    email,
    address
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPersonByID :one
SELECT *
FROM persons
WHERE id = $1;