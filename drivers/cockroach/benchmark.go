package cockroach

import (
	"context"
	"database/sql"

	"github.com/brianvoe/gofakeit"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS data(
    id varchar,
    name varchar,
    value int,
    PRIMARY KEY(id));`

	insertQuery         = `INSERT INTO data (id, name, value) VALUES ($1, $2, $3);`
	selectWhereIdEquals = `SELECT * FROM data WHERE id=$1`
)

type tester struct {
	db *sql.DB
}

func New() *tester {
	db, err := sql.Open("postgres", "postgres://dataset:dataset@localhost:26257/dataset?sslmode=require")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	t := tester{
		db: db,
	}
	t.initTable()

	return &t
}

func (t tester) Write(ctx context.Context) (string, error) {
	id := xid.New().String()
	_, err := t.db.ExecContext(ctx, insertQuery, id, gofakeit.Name(), gofakeit.Number(0, 1000000))
	if err != nil {
		return "", err
	}

	return id, err
}

func (t tester) Read(ctx context.Context, id string) error {
	_, err := t.db.ExecContext(ctx, selectWhereIdEquals, id)

	return err
}

func (t tester) initTable() {
	_, err := t.db.ExecContext(context.Background(), createTable)
	if err != nil {
		panic(err)
	}
}
