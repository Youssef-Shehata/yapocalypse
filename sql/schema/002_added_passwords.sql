-- +goose Up
ALTER TABLE users
ADD COLUMN password Text Not Null ;


-- +goose Down
ALTER TABLE users 
DROP COLUMN password ;
