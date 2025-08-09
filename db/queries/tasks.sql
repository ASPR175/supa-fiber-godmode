-- name: CreateTask :one
INSERT INTO tasks (title,description,user_id)
VALUES ($1,$2,$3)
RETURNING *;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;

-- name: GetTasksByUserId :many
SELECT * FROM tasks
WHERE user_id = $1
ORDER BY created_at DESC;