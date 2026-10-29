-- +goose Up
-- used for 'acting as role'
CREATE TABLE IF NOT EXISTS user_role_assignments (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  organisation_id UUID REFERENCES organisation (id) ON DELETE RESTRICT,
  granted_by UUID NULL REFERENCES users(id) ON DELETE SET NULL,
  expires_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

  CHECK(expires_at IS NULL OR expires_at > created_at)
);

CREATE UNIQUE INDEX uq_user_global_role 
  ON user_role_assignments (user_id,role_id)
  WHERE organisation_id IS NULL AND revoked_at IS NULL;

CREATE UNIQUE INDEX uq_user_organisation_role
  ON user_role_assignments (user_id,role_id,organisation_id)
  WHERE organisation_id IS NOT NULL AND revoked_at IS NULL;

CREATE INDEX idx_user_role_assignments_user
    ON user_role_assignments (user_id);

CREATE INDEX idx_user_role_assignments_role
    ON user_role_assignments (role_id);
-- +goose Down
DROP INDEX IF EXISTS uq_user_global_role;
DROP INDEX IF EXISTS uq_user_organisation_role;
DROP INDEX IF EXISTS idx_user_role_assignments_user;
DROP INDEX IF EXISTS idx_user_role_assignments_role;

DROP TABLE IF EXISTS user_role_assignments;