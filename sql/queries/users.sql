-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, password , username)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3
)
returning *;


-- name: ResetUser :exec
delete  from  users ;

-- name: GetUserById :one
SELECT * from users where id= $1;

-- name: GetUserByUsername :one
SELECT * from users where username= $1;

-- name: GetUserByEmail :one
SELECT * from users where email= $1;

-- name: UpdateUser :one
UPDATE users set email =$1 ,password = $2 , updated_at = now() WHERE id = $3 returning *;


-- name: SubscribeToPremuim :exec 
UPDATE users set premium = true where id = $1;
