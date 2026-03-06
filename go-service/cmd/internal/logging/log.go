package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jackc/pgx/v5"
)

func LogUserAction(ctx context.Context, ch clickhouse.Conn, pg pgx.Tx, action_type string, user_id int) {

	var action_id int
	err := pg.QueryRow(ctx, "SELECT id FROM dicts.actions WHERE label = $1", action_type).Scan(&action_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	query := "INSERT INTO user_actions (dt, user_id, action_id) VALUES ($1, $2, $3)"

	batch, err := ch.PrepareBatch(ctx, query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = batch.Append(time.Now(), user_id, action_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = batch.Send()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
