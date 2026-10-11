-- name: CreateStatementAction :one
INSERT INTO statement_actions (
    statement_id,
    action
)
VALUES (
           $1,
           $2
       )
    RETURNING *;

-- name: ListStatementActions :many
SELECT *
FROM statement_actions
WHERE statement_id = $1
ORDER BY action;

-- name: DeleteStatementAction :exec
DELETE
FROM statement_actions
WHERE statement_id = $1
  AND action = $2;

-- name: DeleteStatementActions :exec
DELETE
FROM statement_actions
WHERE statement_id = $1;

-- name: ListStatementActionsByPolicy :many
select sa.*
from statement_actions sa
join statement s on s.id = sa.statement_id
where s.policy_id = $1
order by sa.statement_id;