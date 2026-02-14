--.sql
CREATE TABLE IF NOT EXISTS notification_groups (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- 002_create_contacts.sql
CREATE TABLE IF NOT EXISTS contacts (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    group_id BIGINT NOT NULL REFERENCES notification_groups(id) ON DELETE CASCADE
);

-- 003_create_notifications.sql
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    group_id BIGINT NOT NULL REFERENCES notification_groups(id) ON DELETE CASCADE
);
