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

	// Find the user
	user, err := queries.GetUserByEmail(ctx, "daniellupo30@gmail.com")
	if err != nil {
		log.Fatalf("user not found: %v", err)
	}

	// Find the organization
	org, err := queries.GetOrganizationByName(ctx, "Acme Corp")
	if err != nil {
		log.Fatalf("org not found: %v", err)
	}

	// Link user to org
	_, err = queries.AddUserToOrganization(ctx, db.AddUserToOrganizationParams{
		UserID:         user.ID,
		OrganizationID: org.ID,
		Role:           db.RoleMember,
	})
	if err != nil {
		fmt.Printf("User already in org: %v\n", err)
	} else {
		fmt.Println("Added user to organization.")
	}

	// Find the project
	project, err := queries.GetProjectByName(ctx, db.GetProjectByNameParams{
		Name:           "Core Engine",
		OrganizationID: org.ID,
	})
	if err != nil {
		log.Fatalf("project not found: %v", err)
	}

	// Link user to project
	_, err = queries.AddUserToProject(ctx, db.AddUserToProjectParams{
		UserID:    user.ID,
		ProjectID: project.ID,
		Role:      db.RoleMember,
	})
	if err != nil {
		fmt.Printf("User already in project: %v\n", err)
	} else {
		fmt.Println("Added user to project.")
	}

	fmt.Println("Done linking user!")
}
