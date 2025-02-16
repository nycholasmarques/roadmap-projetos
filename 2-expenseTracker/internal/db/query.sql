-- name: GetExpense :one
SELECT * FROM expense
WHERE id = $1 LIMIT 1;

-- name: ListExpense :many
SELECT * FROM expense
ORDER BY amount ;

-- name: CreateExpense :one
INSERT INTO expense (
  description, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expense
WHERE id = $1;