-- name: GetPolicy :one
SELECT *
FROM policy
WHERE id = $1;

-- name: ExistsPolicy :one
SELECT EXISTS (
    SELECT 1
    FROM policy
    WHERE resource_type = $1
      AND resource_id = $2
);

-- name: CreatePolicy :one
INSERT INTO policy (resource_type,
                              resource_id,
                              version)
VALUES ($1,
        $2,
        $3) RETURNING *;

-- name: DeletePolicy :exec
DELETE
FROM policy
WHERE id = $1;

-- name: DeletePolicyByResource :exec
DELETE
FROM policy
WHERE resource_type = $1 and resource_id = $2;

-- name: UpdatePolicy :one
UPDATE policy
SET
    version = $2,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: GetPolicyByResource :one
SELECT *
FROM policy
WHERE resource_type = $1
  AND resource_id = $2;

-- name: ListPoliciesByResourceType :many
SELECT *
FROM policy
WHERE resource_type = $1
ORDER BY resource_id;

-- name: ListPolicies :many
SELECT id, resource_type, resource_id, version, created_at, updated_at
FROM policy
WHERE (@resource_type::text = '' OR resource_type = @resource_type::text)
  AND (@resource_id::text = '' OR resource_id = @resource_id::text)
ORDER BY created_at DESC
LIMIT @lim
OFFSET @off;

-- name: CountPolicies :one
SELECT COUNT(*)
FROM policy
WHERE (@resource_type::text = '' OR resource_type = @resource_type::text)
  AND (@resource_id::text = '' OR resource_id = @resource_id::text);