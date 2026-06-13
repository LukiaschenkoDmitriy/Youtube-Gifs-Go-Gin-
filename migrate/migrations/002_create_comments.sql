-- +goose Up
CREATE TABLE comments (
   id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id    UUID NOT NULL,
   video_id   VARCHAR(255) NOT NULL,
   gif_url    VARCHAR(255) DEFAULT NULL,
   text       TEXT DEFAULT NULL,
   answer_to  UUID DEFAULT NULL,
   position     BIGINT GENERATED ALWAYS AS IDENTITY,
   created_at TIMESTAMPTZ DEFAULT NOW(),
   FOREIGN KEY (user_id) REFERENCES users(id),
   FOREIGN KEY (answer_to) REFERENCES comments(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE comments;