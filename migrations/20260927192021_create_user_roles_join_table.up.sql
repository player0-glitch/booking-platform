CREATE TABLE IF NOT EXISTS user_roles (
  user_id INTEGER,
  role_id INTEGER,
  PRIMARY KEY (user_id,role_id)
  CONSTRAINTS fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,,
  CONSTRAINTS fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,,
);
