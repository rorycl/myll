package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestMModel(t *testing.T) {
	dsn := os.Getenv("PG_TEST_DB")
	if dsn == "" {
		t.Fatal("PG_TEST_DB dsn environmental variable not found")
	}
	ctx := context.Background()
	m, err := NewModel(dsn)
	if err != nil {
		t.Fatal(err)
	}

	search_path := []string{"'myll'", `public`}

	user, err := Row[User](
		ctx,
		m,
		search_path,
		"select * from users where id = $1",
		1,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", user)

	type x struct{ A int }
	n, err := Row[x](ctx, m, nil, "select $1::int as a", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", n)

	type y struct{ B string }
	na := pgx.NamedArgs{"b": "hi"}
	o, err := Row[y](ctx, m, nil, "select @b as b", na)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", o)
}

func TestModel(t *testing.T) {
	dsn := os.Getenv("PG_TEST_DB")
	if dsn == "" {
		t.Fatal("PG_TEST_DB dsn environmental variable not found")
	}

	ctx := context.Background()
	m, err := NewModel(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.Exec(ctx, "set search_path=myll"); err != nil {
		t.Fatal(err)
	}
	rows, err := m.Query(ctx, "select * from users where id = $1", 1)
	if err != nil {
		t.Fatal(err)
	}
	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", result)

	if got, want := m.inited, true; got != want {
		t.Errorf("inited got %t want %t", got, want)
	}
	statA := m.Stat().AcquireCount()

	// check no new connection is made
	m, err = NewModel(dsn)
	statB := m.Stat().AcquireCount()

	if got, want := statA, statB; got != want {
		t.Errorf("acquire count got %v want %v", got, want)
	}

	if _, err = m.Exec(ctx, "set search_path=myll"); err != nil {
		t.Fatalf("Exec error: %v\n", err)
	}

	rows, err = m.Query(ctx, "select * from fn_user_manage($1, $2, $3, $4, $5)", "update", "biggles2", "biggles@cl.net", "$2a$06$bgkEkfcteprZC07Djdyx6u/Pw92Z5mWfzqMlUJBwxWcvs.Y/clgQW", "2")
	if err != nil {
		t.Errorf("Query error: %v\n", err)
	}

	result, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", result)

	if _, err = m.Exec(ctx, "set search_path=myll"); err != nil {
		t.Fatalf("Exec error: %v\n", err)
	}

	rows, err = m.Query(ctx, "select * from users where id = $1", 2)
	if err != nil {
		t.Fatal(err)
	}
	result, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", result)

}
