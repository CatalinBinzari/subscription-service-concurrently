package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	// connect to database
	db := initDb()
	db.Ping()

	// create sessios
	session := initSession()

	// create channels

	//  wg

	// set the app config

	// listen to web cxns

}

func initDb() *sql.DB {
	conn := connectToDb()
	if conn == nil {
		log.Panic("Could not connect to the database")
	}
	// connect to db
	return conn
}

func connectToDb() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Printf("Could not connect to the database [%s]. Retrying...", err)
		} else {
			log.Println("Connected to the database")
			return connection
		}

		if counts > 5 {
			log.Panic("Could not connect to the database")
			return nil
		}

		counts++
		log.Print("Backing off for 1 seconds")
		time.Sleep(1 * time.Second)

		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {
	session := scs.New()
	session.Store = redisstore.New(initReddis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session

}

func initReddis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return redisPool
}
