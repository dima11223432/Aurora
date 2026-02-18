CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,  
    username TEXT,                      
    first_name TEXT NOT NULL,          
    last_name TEXT,                   
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE  
);

CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users (telegram_id);


CREATE TABLE IF NOT EXISTS apps (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
