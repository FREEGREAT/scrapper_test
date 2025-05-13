package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron"
)

type Subscription struct {
	ID        int
	ExpiresAt time.Time
	IsValid   bool
}

func main() {
	c := cron.New()

	c.AddFunc("0 0 * * * ", App)

	c.Run()

}

func App() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	expiredSubs, err := getExpSubs(db)
	if err != nil {
		log.Printf("Error getting expired subscriptions: %v", err)
	}

	log.Printf("Found %d expired subscriptions", len(expiredSubs))

	for _, sub := range expiredSubs {
		err := updateSubscriptionStatus(db, sub.ID)
		if err != nil {
			log.Printf("Error updating subscription %d: %v", sub.ID, err)
		} else {
			log.Printf("Subscription %d is expired", sub.ID)
		}
	}

}

func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5430, "user", "user", "subscribe_db")
	return sql.Open("postgres", psqlInfo)
}

func getExpSubs(db *sql.DB) ([]Subscription, error) {
	rows, err := db.Query("SELECT id, expires_at, is_valid FROM subscribers WHERE expires_at < NOW() AND is_valid = true ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subscriptions []Subscription
	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(&sub.ID, &sub.ExpiresAt, &sub.IsValid); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func updateSubscriptionStatus(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE subscribers SET is_valid = false WHERE id = $1", id)
	return err
}
