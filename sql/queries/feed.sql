-- name: GetInitialFeed :many 
select y.* from yaps y join feed f on y.user_id = f.user_id
where y.user_id = $1 order by y.created_at Desc LIMIT 20 ;


