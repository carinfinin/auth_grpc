CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    pass_hash BLOB NOT NULL

);
CREATE INDEX IF NOT EXISTS idx_email ON users (email)