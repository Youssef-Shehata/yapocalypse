

-- name: GetFollowersOf :many
SELECT u.id , u.email,u.created_at , u.updated_at , u.username  ,u.premuim from users u join followers f on u.id = f.follower_id where f.followee_id == $1 order by f.created_at DESC;



-- name: GetFolloweesOf :many
SELECT u.id , u.email,u.created_at , u.updated_at , u.username  ,u.premuim from users u join followers f on u.id = f.followee_id where f.follower_id == $1 order by f.created_at DESC;


-- name: AddFollower :exec
iNSERT INTO followers (follower_id, followee_id)
VALUES ($1, $2);




