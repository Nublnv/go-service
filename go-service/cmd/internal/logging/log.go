package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActionData struct {
	actionType string
	userid     int
}

func LogUserAction(ctx context.Context, actionType string, userid int) {
	var c chan<- *ActionData = ctx.Value("logging").(chan *ActionData)

	c <- &ActionData{
		actionType: actionType,
		userid:     userid,
	}
}

func DoLoggingActions(ctx context.Context, ch clickhouse.Conn, pg *pgxpool.Conn) {
	var c <-chan *ActionData = ctx.Value("logging").(chan *ActionData)
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c:
			var action_id int
			err := pg.QueryRow(ctx, "SELECT id FROM dicts.actions WHERE label = $1", data.actionType).Scan(&action_id)
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
			err = batch.Append(time.Now(), data.userid, action_id)
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
	}
}
