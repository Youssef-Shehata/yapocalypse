-- name: NewYap :one

INSERT INTO Yaps(id, created_at, updated_at, user_id ,body )

VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: ResetYaps :exec
DELETE from  Yaps;


-- name: GetYapsByUserId :many
select * from Yaps where user_id = $1 order by created_at desc ;
-- name: GetYapById :one
select * from Yaps where id = $1;

-- name: DeleteYap :exec
DELETE from  Yaps where id = $1 and user_id = $2;


