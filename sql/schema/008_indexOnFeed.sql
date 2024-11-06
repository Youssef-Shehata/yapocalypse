
-- +goose Up 
CREATE INDEX feed_idx ON Feed(user_id);

-- +goose Down
DROP INDEX if exists feed_idx ON Feed(user_id);
