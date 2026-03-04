package migrations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoPgMigrates(ctx context.Context, conn pgxpool.Conn) error {
	path := os.Getenv("MIGRATIONS_PATH")
	files, err := os.ReadDir(fmt.Sprintf("%s/postgres", path))
	if err != nil {
		return err
	}

	migrationsMapping, err := getMigrations(ctx, conn)
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

		query := "INSERT INTO migrations.postgres (migration_date, name, application_date) VALUES ($1, $2, $3)"

		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, string(fileData))
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		_, err = tx.Exec(ctx, query, date, name, time.Now())
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		tx.Commit(ctx)
	}
	return nil
}

func parseFileName(filename string) (*time.Time, string, error) {
	re := regexp.MustCompile(`(\d{14})_(\w+).sql`)

	result := re.FindStringSubmatch(filename)
	if len(result) == 3 {
		date, err := time.Parse("20060102150405", result[1])
		if err != nil {
			return nil, "", errors.New("Wrong time format")
		}
		name := strings.Replace(result[2], "_", " ", 0)
		return &date, name, nil
	} else {
		return nil, "", fmt.Errorf("Wrong file %s name format", filename)
	}

}

func getMigrations(ctx context.Context, conn pgxpool.Conn) (map[*time.Time]string, error) {
	var migrationsMap = map[*time.Time]string{}

	rows, err := conn.Query(ctx, "SELECT migration_date, name FROM migrations.postgres")
	if err != pgx.ErrNoRows && err != nil {
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
