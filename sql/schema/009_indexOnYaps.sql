

-- +goose Up 
CREATE INDEX yap_idx ON Yaps(user_id);

-- +goose Down
DROP INDEX if exists yap_idx ON Yaps(user_id);
