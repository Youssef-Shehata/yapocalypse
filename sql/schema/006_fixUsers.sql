-- +goose Up
Alter table users 
Add COLUMN username text unique not null ;



-- +goose Down
alter table users 
drop COLUMN username ;




