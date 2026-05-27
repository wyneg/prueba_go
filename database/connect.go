package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5" //<=== driver necesario para la conexión
	"github.com/joho/godotenv"
)

func Connect() (*pgx.Conn, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())

	fmt.Println("Conexión da base de datos PostgreSQL exitosa")

	return conn, nil
}
