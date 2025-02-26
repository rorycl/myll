package main

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64 `db:"id"`
	Created   time.Time
	Modified  time.Time
	Name      string
	Email     string
	HPassword string `db:"hpassword"`
}

func dbquery(url string, searchpath string, user int) (*User, error) {

	ctx := context.Background()
	db, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("Connection error : %s : %v\n", url, err)
	}
	_, err = db.Exec(ctx, "set search_path="+searchpath)
	if err != nil {
		return nil, fmt.Errorf("Exec error: %v\n", err)
	}

	var results []*User
	err = pgxscan.Select(ctx, db, &results, "select * from users where id = 1")
	if err != nil {
		return nil, fmt.Errorf("Exec error: %v\n", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("expected 1 result, got %d", len(results))
	}

	return results[0], nil
}

func dbquery2(url string, searchpath string, user int) (*User, error) {

	ctx := context.Background()
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("Connection error : %s : %v\n", url, err)
	}
	defer db.Close()

	if _, err = db.Exec(ctx, "set search_path="+searchpath); err != nil {
		return nil, fmt.Errorf("Exec error: %v\n", err)
	}

	rows, err := db.Query(ctx, "select * from users where id = $1", user)
	if err != nil {
		return nil, fmt.Errorf("failed query: %w", err)
	}
	// rows closed by CollectExactlyOneRow

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, fmt.Errorf("collect error: %v\n", err)
	}

	return &result, nil
}

func main() {
	// ctx := context.Background()
	url := "postgres://mylluser:Oof7eiyuv:@127.0.0.1:5432/myll"
	/*
		db, err := pgxpool.New(ctx, url)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var users []*User
		err = pgxscan.Select(ctx, db, &users, `select * from myll.users where id = 1`)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("#v\n", users)
	*/
	u, err := dbquery(url, "myll", 1)
	fmt.Printf("%#v %v\n\n", u, err)
	u, err = dbquery2(url, "myll", 1)
	fmt.Printf("%#v %v\n\n", u, err)
}
