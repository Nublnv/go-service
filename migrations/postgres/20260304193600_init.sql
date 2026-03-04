-- DROP SCHEMA auth;

CREATE SCHEMA auth AUTHORIZATION postgres;


-- Permissions

GRANT ALL ON SCHEMA auth TO postgres;

-- auth.users определение

-- Drop table

-- DROP TABLE auth.users;

CREATE TABLE auth.users (
	login varchar NOT NULL,
	passhash bytea NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (login)
);

-- Permissions

ALTER TABLE auth.users OWNER TO postgres;
GRANT ALL ON TABLE auth.users TO postgres;

-- DROP SCHEMA migrations;

CREATE SCHEMA migrations AUTHORIZATION postgres;


-- Permissions

GRANT ALL ON SCHEMA migrations TO postgres;


CREATE TABLE migrations.postgres (
	migration_date timestamp with time zone NOT NULL,
	"name" varchar NOT NULL,
	application_date timestamp with time zone NOT NULL,
	CONSTRAINT postgres_pk PRIMARY KEY (migration_date,"name")
);
