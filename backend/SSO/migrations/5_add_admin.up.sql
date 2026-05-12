INSERT INTO users(user_id, telegram_id, username, first_name, last_name, is_admin)
VALUES (1, 123456789, 'john_doe', 'john', 'doe', true)
ON CONFLICT (user_id) 
DO UPDATE SET 
    username = EXCLUDED.username,
    first_name = EXCLUDED.first_name,
    is_admin = EXCLUDED.is_admin;
