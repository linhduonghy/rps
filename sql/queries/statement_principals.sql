-- name: CreateStatementPrincipal :one
INSERT INTO statement_principals (
    statement_id,
    principal_type,
    principal_id
)
VALUES (
           $1,
           $2,
           $3
       )
    RETURNING *;

-- name: ListStatementPrincipals :many
SELECT *
FROM statement_principals
WHERE statement_id = $1
ORDER BY principal_type, principal_id;

-- name: DeleteStatementPrincipal :exec
DELETE
FROM statement_principals
WHERE statement_id = $1
  AND principal_type = $2
  AND principal_id = $3;

-- name: DeleteStatementPrincipals :exec
DELETE
FROM statement_principals
WHERE statement_id = $1;

-- name: ListStatementPrincipalsByPolicy :many
select sp.*
from statement_principals sp
join statement s on s.id = sp.statement_id
where s.policy_id = $1
order by sp.statement_id;