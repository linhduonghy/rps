-- name: CreateStatementResource :one
INSERT INTO statement_resources (
    statement_id,
    resource
)
VALUES (
           $1,
           $2
       )
    RETURNING *;

-- name: ListStatementResources :many
SELECT *
FROM statement_resources
WHERE statement_id = $1
ORDER BY resource;

-- name: DeleteStatementResource :exec
DELETE
FROM statement_resources
WHERE statement_id = $1
  AND resource = $2;

-- name: DeleteStatementResources :exec
DELETE
FROM statement_resources
WHERE statement_id = $1;

-- name: ListStatementResourcesByPolicy :many
select sr.*
from statement_resources sr
         join statement s on s.id = sr.statement_id
where s.policy_id = $1
order by sr.statement_id;