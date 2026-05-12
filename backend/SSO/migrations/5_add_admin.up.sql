INSERT into users(user_id, telegram_id, username, first_name, last_name, is_admin)
VALUES (1, 123456789, 'john_doe', 'john', 'doe', true)
ON CONFLICT DO NOTHING;
