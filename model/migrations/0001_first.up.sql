CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(64),
  password_hash VARCHAR(64),
  email VARCHAR(128),
  created_at INT64,
  updated_at INT64
);

CREATE TABLE sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  token VARCHAR(64),
  user_id INT64,
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
  content_type VARCHAR(256),
  uuid VARCHAR(36),
  note_id INT64,
  created_at INT64,
  updated_at INT64
);
