CREATE SCHEMA dicts AUTHORIZATION postgres;

GRANT ALL ON SCHEMA dicts TO postgres;

CREATE TABLE dicts.actions (
	id Int8 NOT NULL,
    "label" VARCHAR NOT NULL,
    CONSTRAINT actions_pk PRIMARY KEY (id)
);

INSERT INTO dicts.actions (id, "label") VALUES
( 1000, 'registration' ),
( 1001, 'login' );