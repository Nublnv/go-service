package migrations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Nublnv/go-service/cmd/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration = struct {
	date *time.Time
	name string
}

func DoPgMigrates(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Conn().Close(ctx)
	path := os.Getenv("MIGRATIONS_PATH")
	files, err := os.ReadDir(fmt.Sprintf("%s/postgres", path))
	if err != nil {
		return err
	}

	migrationsMapping, err := getMigrations(ctx, conn, "postgres")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		fileData, err := os.ReadFile(fmt.Sprintf("%s/postgres/%s", path, filename))
		if err != nil {
			return err
		}
		date, name, err := parseFileName(filename)
		if err != nil {
			return err
		}
		value, ok := migrationsMapping[date]
		if ok && value == name {
			continue
		}

		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, string(fileData))
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		tx.Commit(ctx)
		err = addMigrationRecord(ctx, conn, &Migration{
			date: date,
			name: name,
		}, "postgres")
		if err != nil {
			return err
		}
	}
	return nil
}

func DoChickMigrations(ctx context.Context, pool *db.Pool, pgPool *pgxpool.Pool) error {
	pgConn, err := pgPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pgConn.Conn().Close(ctx)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	path := os.Getenv("MIGRATIONS_PATH")
	files, err := os.ReadDir(fmt.Sprintf("%s/clickhouse", path))
	if err != nil {
		return err
	}

	migrationsMapping, err := getMigrations(ctx, pgConn, "clickhouse")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		fileData, err := os.ReadFile(fmt.Sprintf("%s/clickhouse/%s", path, filename))
		if err != nil {
			return err
		}
		date, name, err := parseFileName(filename)
		if err != nil {
			return err
		}
		value, ok := migrationsMapping[date]
		if ok && value == name {
			continue
		}

		err = conn.Conn().Exec(ctx, string(fileData))
		if err != nil {
			return err
		}
		err = addMigrationRecord(ctx, pgConn, &Migration{
			date: date,
			name: name,
		}, "clickhouse")
		if err != nil {
			return err
		}
	}
	return nil
}

func addMigrationRecord(ctx context.Context, conn *pgxpool.Conn, migration *Migration, table string) error {
	query := "INSERT INTO migrations.$4 (migration_date, name, application_date) VALUES ($1, $2, $3)"

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	if migration.date == nil {
		migration.date = new(time.Time)
		*migration.date = time.Unix(0, 0)
	}
	_, err = tx.Exec(ctx, query, &migration.date, migration.name, time.Now(), table)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	tx.Commit(ctx)
	return nil
}

func parseFileName(filename string) (*time.Time, string, error) {
	re := regexp.MustCompile(`(\d{14})_(\w+).sql`)

	result := re.FindStringSubmatch(filename)
	if len(result) == 3 {
		name := strings.Replace(result[2], "_", " ", 0)
		if result[1] != "00000000000000" {
			date, err := time.Parse("20060102150405", result[1])
			if err != nil {
				return nil, "", errors.New("Wrong time format")
			}
			return &date, name, nil
		} else {
			return nil, name, nil
		}
	} else {
		return nil, "", fmt.Errorf("Wrong file %s name format", filename)
	}

}

func getMigrations(ctx context.Context, conn *pgxpool.Conn, table string) (map[*time.Time]string, error) {
	var migrationsMap = map[*time.Time]string{}

	rows, err := conn.Query(ctx, "SELECT migration_date, name FROM migrations.$1", table)
	if err != pgx.ErrNoRows && err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "42P01" {
				return migrationsMap, nil
			}
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var date time.Time
		if err := rows.Scan(&date, name); err != nil {
			return nil, err
		}
		migrationsMap[&date] = name
	}
	return migrationsMap, nil
}
