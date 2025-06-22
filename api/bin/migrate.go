package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var uname string
	var pass string
	var db string
	var host string

	flag.StringVar(&uname, "u", "postgres", "Username for postgresql. Defaults to postgres")
	flag.StringVar(&pass, "p", "mysecretpassword", "Password for postgres user. Defaults to mysecretpassword.")
	flag.StringVar(&db, "d", "postgres", "Database name. Defaults to postgres.")
	flag.StringVar(&host, "h", "localhost", "Database host. Defaults to localhost.")

	dbUrl := "postgres://"+uname+":"+pass+"@"+host+":5432/"+db+"?sslmode=disable"
	m, err := migrate.New(
		"file://db/migrations",
        dbUrl,
    )
    if err != nil {
        log.Fatal(err)
    }
    if err := m.Up(); err != nil {
        log.Fatal(err)
    }
    m.Close()
}