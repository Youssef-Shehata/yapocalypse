
-- +goose Up
ALTER TABLE users
ADD COLUMN premuim bool Not Null DEFAULT false;


-- +goose Down
ALTER TABLE users 
DROP COLUMN premuim;
