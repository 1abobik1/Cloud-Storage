-- Удаление индексов
DROP INDEX IF EXISTS idx_refresh_token_user_id;
DROP INDEX IF EXISTS idx_refresh_token_user_id_token;

-- Удаление таблицы refresh_token
DROP TABLE IF EXISTS refresh_token;

-- Удаление таблицы auth_users
DROP TABLE IF EXISTS auth_users;

-- Удаление пользовательского enum типа данных
DROP TYPE IF EXISTS token_platform_enum;