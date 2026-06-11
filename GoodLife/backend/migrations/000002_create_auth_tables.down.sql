-- Rollback auth tables
DROP TABLE IF EXISTS auth.sessions CASCADE;
DROP TABLE IF EXISTS auth.user_roles CASCADE;
DROP TABLE IF EXISTS auth.roles CASCADE;
DROP TABLE IF EXISTS auth.users CASCADE;
