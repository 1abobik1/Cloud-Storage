CREATE TABLE IF NOT EXISTS auth_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(130) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_activated BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS refresh_token (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    user_id INT NOT NULL,
    ip INET NOT NULL,
    FOREIGN KEY (user_id) REFERENCES auth_users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_token_user_id ON  refresh_token (user_id);
CREATE INDEX idx_refresh_token_ip ON refresh_token (ip);
CREATE INDEX idx_refresh_token_user_id_token ON refresh_token (user_id, token);
