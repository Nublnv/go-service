CREATE TABLE migrations.clickhouse (
	migration_date timestamp with time zone NOT NULL,
	"name" varchar NOT NULL,
	application_date timestamp with time zone NOT NULL,
	CONSTRAINT clickhouse_pk PRIMARY KEY (migration_date,"name")
);