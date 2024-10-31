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
