-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, password)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
returning *;


-- name: ResetUser :exec
delete  from  users ;

-- name: GetUserByEmail :one
SELECT * from users where email= $1;

