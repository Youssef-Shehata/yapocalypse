
-- +goose Up
CREATE TABLE Users (
id uuid PRIMARY KEY,
created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
email Text unique not null
);

CREATE TABLE Tweets(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
body TEXT NOT NULL,
user_id UUID NOT NULL,
 CONSTRAINT fk_user_id
      FOREIGN KEY(user_id)
        REFERENCES users(id )
        ON DELETE CASCADE 
);

-- +goose Down
DROP  TABLE Tweets;
DROP TABLE users;
