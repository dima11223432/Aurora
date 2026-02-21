CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,  
    username TEXT,                      
    first_name TEXT NOT NULL,          
    last_name TEXT,                   
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users (telegram_id);


CREATE TABLE IF NOT EXISTS channels (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    username TEXT NOT NULL,

    CONSTRAINT fk_channels_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT uq_channels_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS apps (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
