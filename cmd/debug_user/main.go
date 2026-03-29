package main

import (
	"context"
	"fmt"
	"log"

	"github.com/luponetn/enx/internal/config"
	"github.com/luponetn/enx/internal/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := db.ConnectDB(cfg.DbUrl)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	ctx := context.Background()

	user, err := queries.GetUserByEmail(ctx, "daniellupo30@gmail.com")
	if err != nil {
		fmt.Printf("User 'daniellupo30@gmail.com' not found: %v\n", err)
		return
	}
	fmt.Printf("User: %s (ID: %v)\n", user.Email, user.ID)

	projects, err := queries.GetProjectsByUser(ctx, user.ID)
	if err != nil {
		fmt.Printf("Could not get projects: %v\n", err)
		return
	}
	fmt.Printf("Found %d projects for this user.\n", len(projects))
	for _, p := range projects {
		fmt.Printf("- %s (ID: %v) in Org: %v\n", p.Name, p.ID, p.OrganizationID)
	}
}
