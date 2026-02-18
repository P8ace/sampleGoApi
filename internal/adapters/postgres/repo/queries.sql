-- name: ListProducts :many
SELECT *
FROM products;

-- name: FindProductsById :one
SELECT *
FROM products
WHERE id = $1;
