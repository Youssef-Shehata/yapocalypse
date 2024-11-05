
-- +goose Up 
CREATE TABLE Followers(
    follower_id uuid references Users(id),
    followee_id uuid references Users(id),
    created_at timestamp Not null Default now(),
    Primary Key (follower_id,followee_id)
);




-- +goose Down
DROP TABLE Followers;
