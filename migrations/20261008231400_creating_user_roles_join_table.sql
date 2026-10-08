-- +goose Up
CREATE TABLE IF NOT EXISTS user_roles (
  user_id INTEGER,
  role_id INTEGER,
  PRIMARY KEY (user_id,role_id),
  CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- +goose Down
-- Disabling so that a clean drop order can be maintained
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS user_roles;

PRAGMA foreign_keys = ON;