CREATE TABLE users (
                      id SERIAL PRIMARY KEY,
                      username VARCHAR(255) UNIQUE NOT NULL,
                      password VARCHAR(255) NOT NULL
);

-- DO $$
-- BEGIN
--     IF EXISTS (
--         SELECT 1
--         FROM pg_constraint
--         WHERE conname = 'uni_users_username'
--     ) THEN
-- ALTER TABLE users DROP CONSTRAINT uni_users_username;
-- END IF;
-- END $$;