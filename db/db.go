package db

import (
	"context"
	"log"
	"todoGoApi/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB_CONN *pgxpool.Pool

func InitDatabase(){
	var DB_ERR error
	DB_CONN, DB_ERR = pgxpool.New(context.Background(), utils.DB_URL)
	if DB_ERR != nil {
		log.Fatalln("Unable to connect to database: ", DB_ERR)
	}
	log.Println("Database connection established")
}

func CloseDatabase(){
	DB_CONN.Close()
	log.Println("Database connection disconnected")
}