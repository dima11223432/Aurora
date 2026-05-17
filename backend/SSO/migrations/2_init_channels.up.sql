
CREATE TABLE IF NOT EXISTS channels (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    channel_username TEXT NOT NULL,

    CONSTRAINT fk_channels_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT uq_channels_username UNIQUE (channel_username)
);
