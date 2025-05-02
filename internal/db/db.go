package db

import (
	"context"
	"fmt"
	"github.com/faceoffdev/edutracking-bot/internal/config"
	"github.com/jackc/pgx/v5"
	"iter"
)

var conn *pgx.Conn

func Init(c config.PostgresConfig) error {
	connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&application_name=edutracking_bot", c.User, c.Pass, c.Host, c.Port, c.Dbname)

	var err error
	conn, err = pgx.Connect(context.Background(), connstr)

	return err
}

func doQuery[T any](query string, args ...any) (iter.Seq2[T, error], error) {
	rows, err := conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	return func(yield func(T, error) bool) {
		defer rows.Close()

		for rows.Next() {
			value, err := pgx.RowToStructByName[T](rows)
			if !yield(value, err) {
				return
			}
		}
	}, nil
}

func getList[T any](rows iter.Seq2[T, error]) []T {
	var result []T
	for row, err := range rows {
		if err == nil {
			result = append(result, row)
		}
	}

	return result
}
