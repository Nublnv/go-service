CREATE INDEX users_user_id_idx ON auth.users (user_id);
ALTER TABLE auth.users ADD email varchar NOT NULL DEFAULT '';
ALTER TABLE auth.users ADD CONSTRAINT users_unique UNIQUE (email);
