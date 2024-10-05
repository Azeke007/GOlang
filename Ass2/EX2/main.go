package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name string
	Age  int
}

func connectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=AK_qwerty dbname=GOlang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to PostgreSQL using GORM")
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Table created successfully with AutoMigrate")

	users := []User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	for _, user := range users {
		db.Create(&user)
		fmt.Printf("Inserted user %s, age %d\n", user.Name, user.Age)
	}

	var allUsers []User
	db.Find(&allUsers)

	fmt.Println("Querying all users:")
	for _, user := range allUsers {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}