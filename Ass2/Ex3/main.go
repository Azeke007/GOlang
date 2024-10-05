package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var dbSQL *sql.DB
var dbGORM *gorm.DB

func connectSQL() *sql.DB {
	connStr := "user=postgres password=AK_qwerty dbname=GOlang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return db
}

func connectGORM() *gorm.DB {
	dsn := "host=localhost user=postgres password=AK_qwerty dbname=GOlang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	db.AutoMigrate(&User{})
	return db
}

func getUsersSQL(w http.ResponseWriter, r *http.Request) {
	rows, err := dbSQL.Query("SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func createUserSQL(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	_, err := dbSQL.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getUsersGORM(w http.ResponseWriter, r *http.Request) {
	var users []User
	dbGORM.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func createUserGORM(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	dbGORM.Create(&user)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	dbSQL = connectSQL()
	dbGORM = connectGORM()

	r := mux.NewRouter()

	r.HandleFunc("/sql/users", getUsersSQL).Methods("GET")
	r.HandleFunc("/sql/users", createUserSQL).Methods("POST")

	r.HandleFunc("/gorm/users", getUsersGORM).Methods("GET")
	r.HandleFunc("/gorm/users", createUserGORM).Methods("POST")

	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
