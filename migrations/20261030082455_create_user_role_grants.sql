-- +goose Up
-- If organisation already exists, user their actual table name and type
CREATE TABLE IF NOT EXISTS user_role_grants (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL 
   REFERENCES users(id) ON DELETE CASCADE,
  organisation_id UUID NULL 
    REFERENCES organisations (id) ON DELETE CASCADE,
  role_id INTEGER NOT NULL 
    REFERENCES roles (id) ON DELETE RESTRICT,
  granted_by UUID NOT NULL
   REFERENCES users(id) ON DELETE SET NULL,
  expires_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_role_grants_user_id
  ON user_role_grants(user_id);

CREATE INDEX IF NOT EXISTS idx_user_role_grants_role_id
  ON user_role_grants(role_id);
-- +goose Down
DROP INDEX IF EXISTS idx_user_role_grants_role_id;
DROP INDEX IF EXISTS idx_user_role_grants_user_id
DROP TABLE IF NOT EXISTS user_role_grants;

