
-- +goose Up
CREATE INDEX idx_followers_followee_id ON Followers(followee_id);

-- +goose Down
Drop INDEX IF EXISTS  idx_followers_followee_id;
