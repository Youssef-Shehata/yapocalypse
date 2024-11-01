-- name: CreateTweet :one
INSERT INTO Tweets(id, created_at, updated_at, user_id ,body )
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: ResetTweets :exec
DELETE from  Tweets;

-- name: GetTweets :many
select * from Tweets where user_id = $1 order by created_at;

-- name: GetTweetById :one
select * from Tweets where id = $1;
