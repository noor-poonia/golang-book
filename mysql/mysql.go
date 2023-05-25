package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func getMYSQLConnectionString() string {
	var connect string
	// host := os.Getenv("MYSQL_HOST")
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DBNAME")

	connect = fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", username, password, port, dbName)
	return connect
}

func MysqlConnction() {
	// db , err := sql.Open("mysql", getMYSQLConnectionString())
	dbConn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books")
    defer dbConn.Close()
	if err!= nil {
        log.Fatal(err)
    }

	err = dbConn.Ping()
	if err!= nil {
        log.Fatal(err)
    }
	fmt.Println("Successfully connected to mysql")

	db = dbConn

	// err = createBooksTable()
	// if err!= nil {
    //     log.Fatal(err)
    // }

	// err = consumeMessage()
	// if err!= nil {
    //     log.Fatal(err)
    // }

	// return dbConn
}

func MysqlTable(db *sql.DB) error {	
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS book (
		id INT AUTO_INCREMENT PRIMARY KEY, 
		name VARCHAR(255) NOT NULL, 
		author VARCHAR(255) NOT NULL, 
		pages INT
	)`)
	
	if err!= nil {
        return err
    }

	return nil
}

// one function will cal the db connection, then the table creation and then the consume msg function