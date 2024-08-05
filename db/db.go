package db

import (
	"context"
	"fmt"
	"github.com/tantanok221/cloudflare-image-backend/models"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/tantanok221/cloudflare-image-backend/internal/helper"
)

var (
	PASSWORD string = helper.GetEnv("PASSWORD")
	USER     string = helper.GetEnv("USERS")
	HOST     string = helper.GetEnv("HOST")
	DBNAME   string = helper.GetEnv("DBNAME")
	PORT     string = helper.GetEnv("PORT")
)

func Init() *models.Queries {
	dsn := fmt.Sprintf("%v://%v:%v@%v:%v/postgres", DBNAME, USER, PASSWORD, HOST, PORT)
	//dsn := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v", USER, PASSWORD, HOST, PORT, DBNAME)
	print(dsn)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalln("Unable to connect to database", err)
	}
	queries := models.New(conn)
	return queries
}
