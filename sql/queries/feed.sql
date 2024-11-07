-- name: GetFeed :many 
select y.* from yaps y join feed f on y.user_id = f.user_id
where y.user_id = $1 order by y.created_at Desc offset $2 LIMIT 20 ;


-- name: AddToFeed :exec
 INSERT INTO feed (user_id , yap_id)
VALUES($1,$2);


