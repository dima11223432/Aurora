INSERT INTO users(telegram_id, username, first_name, last_name, is_admin)
VALUES (123456789, 'john_doe', 'john', 'doe', true)
ON CONFLICT (telegram_id) 
DO UPDATE SET 
    username = EXCLUDED.username,
    first_name = EXCLUDED.first_name,
    is_admin = EXCLUDED.is_admin;
