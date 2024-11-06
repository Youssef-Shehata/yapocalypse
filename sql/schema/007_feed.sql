-- +goose Up 
CREATE TABLE Feed (
    user_id uuid references users(id) on delete cascade ,
    yap_id  uuid references yaps(id) on delete cascade ,
    primary key (user_id, yap_id)
);

-- +goose Down
DROP TABLE Feed ;
