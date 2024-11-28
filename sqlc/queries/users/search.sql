-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: SearchUsers :many
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE 
    CASE 
        WHEN $1::text != '' THEN
            name ILIKE '%' || $1 || '%' OR
            email ILIKE '%' || $1 || '%'
        ELSE true
    END
ORDER BY created_at DESC
LIMIT $2 OFFSET $3; 