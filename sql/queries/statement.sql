-- name: CreateStatement :one
INSERT INTO statement (
    policy_id,
    sid,
    effect
)
VALUES (
           $1,
           $2,
           $3
       )
    RETURNING *;

-- name: GetStatement :one
SELECT *
FROM statement
WHERE id = $1;

-- name: ListStatementsByPolicy :many
SELECT *
FROM statement
WHERE policy_id = $1
ORDER BY id;

-- name: DeleteStatement :exec
DELETE
FROM statement
WHERE id = $1;

-- name: DeleteStatementsByPolicy :exec
DELETE
FROM statement
WHERE policy_id = $1;