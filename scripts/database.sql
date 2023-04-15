/*
DROP TABLE blogs;
DROP TABLE user_roles;
DROP TABLE users;
DROP TABLE roles;
*/

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS blogs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  tags TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    permissions TEXT[]
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id INT REFERENCES users(id),
  role_id INT REFERENCES roles(id),
  PRIMARY KEY (user_id, role_id)
);

INSERT INTO roles (name, permissions) VALUES
('admin', ARRAY['create_user', 'read_user', 'update_user', 'delete_user', 'create_blog', 'read_blog', 'update_blog', 'delete_blog']),
('editor', ARRAY['create_blog', 'read_blog', 'update_blog', 'delete_blog']),
('viewer', ARRAY['read_user', 'read_blog']);
