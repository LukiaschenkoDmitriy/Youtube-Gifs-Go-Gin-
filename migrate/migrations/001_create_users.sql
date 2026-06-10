-- 001_create_users.sql

-- +goose Up
CREATE TABLE users (
   id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   oauth_id   VARCHAR(255) NOT NULL,
   picture    VARCHAR(255) NOT NULL,
   name       VARCHAR(255) NOT NULL,
   email      VARCHAR(255) UNIQUE NOT NULL,
   custom_url VARCHAR(255) DEFAULT NULL,
   created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;