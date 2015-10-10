CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email VARCHAR(128) UNIQUE,
  password_hash VARCHAR(64),
  created_at INT64,
  updated_at INT64
);

CREATE TABLE sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  token VARCHAR(32) UNIQUE,
  user_id INT64,
  expires_at INT64,
  created_at INT64,
  updated_at INT64
);

CREATE TABLE notes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title VARCHAR(256),
  content TEXT,
  owner_id INT64,
  share_state INT8,
  created_at INT64,
  updated_at INT64
);

CREATE TABLE images (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  content_type VARCHAR(32),
  uuid VARCHAR(36) UNIQUE,
  url VARCHAR(128),
  note_id INT64,
  created_at INT64,
  updated_at INT64
);
