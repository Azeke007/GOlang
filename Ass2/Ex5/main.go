package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name    string   `gorm:"not null"`
	Age     int      `gorm:"not null"`
	Profile Profile  
}

type Profile struct {
	gorm.Model
	UserID           uint   `gorm:"uniqueIndex"`
	Bio              string
	ProfilePictureURL string
}

func connectGORM() *gorm.DB {
	dsn := "host=localhost user=postgres password=AK_qwerty dbname=GOlang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to configure connection pool:", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	return db
}

func main() {
	db := connectGORM()

	err := db.AutoMigrate(&User{}, &Profile{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	createUserWithProfile(db)

	queryUserWithProfile(db)

	updateUserProfile(db, 1, "Updated Bio", "https://updated.picture.url")

	deleteUserWithProfile(db, 1)
}

func createUserWithProfile(db *gorm.DB) {
	user := User{
		Name: "Alice",
		Age:  30,
		Profile: Profile{
			Bio:              "I love hiking and outdoor adventures.",
			ProfilePictureURL: "https://example.com/picture.jpg",
		},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		fmt.Println("User and profile created successfully")
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create user and profile: %v", err)
	}
}

func queryUserWithProfile(db *gorm.DB) {
	var users []User
	db.Preload("Profile").Find(&users)

	fmt.Println("Queried Users with Profiles:")
	for _, user := range users {
		fmt.Printf("User: %s, Age: %d, Bio: %s, Profile Picture: %s\n",
			user.Name, user.Age, user.Profile.Bio, user.Profile.ProfilePictureURL)
	}
}

func updateUserProfile(db *gorm.DB, userID uint, newBio, newProfilePictureURL string) {
	var profile Profile
	db.Where("user_id = ?", userID).First(&profile)
	profile.Bio = newBio
	profile.ProfilePictureURL = newProfilePictureURL

	db.Save(&profile)
	fmt.Println("User profile updated successfully")
}

func deleteUserWithProfile(db *gorm.DB, userID uint) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&Profile{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&User{}, userID).Error; err != nil {
			return err
		}
		fmt.Println("User and associated profile deleted successfully")
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to delete user and profile: %v", err)
	}
}
