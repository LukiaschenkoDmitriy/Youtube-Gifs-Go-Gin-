-- 003_create_comments_action_table.sql

-- +goose Up
CREATE TABLE comment_actions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  comment_id UUID NOT NULL,
  user_id    UUID  NOT NULL,
  type       INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES  users(id),
  FOREIGN KEY (comment_id) REFERENCES  comments(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE comment_actions;