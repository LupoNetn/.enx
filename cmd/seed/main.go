package main

import (
	"context"
	"fmt"
	"log"

	"github.com/luponetn/enx/internal/config"
	"github.com/luponetn/enx/internal/db"
	"github.com/luponetn/enx/internal/utils"
)

func main() {
	// Setup and load app config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	pool, err := db.ConnectDB(cfg.DbUrl)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	ctx := context.Background()

	fmt.Println("Starting database seeding...")

	// 1. Create User
	hashedPassword, _ := utils.HashPassword("password123")
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:    "admin@enx.com",
		Name:     "Admin User",
		Password: hashedPassword,
	})
	
	var userID = user.ID
	if err != nil {
		fmt.Printf("User 'admin@enx.com' might already exist: %v\n", err)
		u, getErr := queries.GetUserByEmail(ctx, "admin@enx.com")
		if getErr != nil {
			log.Fatalf("failed to create or fetch user: %v", getErr)
		}
		userID = u.ID
	} else {
		fmt.Printf("Created user: %s (ID: %v)\n", user.Email, user.ID)
	}

	// 2. Create Organization
	org, err := queries.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:      "Acme Corp",
		Email:     "contact@acme.com",
		Passkey:   "acme-passkey-123",
		CreatedBy: userID,
	})
	
	var orgID = org.ID
	if err != nil {
		fmt.Printf("Organization 'Acme Corp' might already exist: %v\n", err)
		o, getErr := queries.GetOrganizationByName(ctx, "Acme Corp")
		if getErr != nil {
			log.Fatalf("failed to create or fetch organization: %v", getErr)
		}
		orgID = o.ID
	} else {
		fmt.Printf("Created organization: %s (ID: %v)\n", org.Name, org.ID)
	}

	// 3. Add User to Org
	_, err = queries.AddUserToOrganization(ctx, db.AddUserToOrganizationParams{
		UserID:         userID,
		OrganizationID: orgID,
		Role:           db.RoleOwner,
	})
	if err != nil {
		fmt.Printf("Note: User-Org link might already exist: %v\n", err)
	} else {
		fmt.Println("Linked user to organization as Owner.")
	}

	// 4. Create Project
	project, err := queries.CreateProject(ctx, db.CreateProjectParams{
		Name:           "Core Engine",
		Passkey:        "engine-secret-456",
		OrganizationID: orgID,
		CreatedBy:      userID,
	})
	
	var projectID = project.ID
	if err != nil {
		fmt.Printf("Project 'Core Engine' might already exist: %v\n", err)
		p, getErr := queries.GetProjectByName(ctx, db.GetProjectByNameParams{
			Name:           "Core Engine",
			OrganizationID: orgID,
		})
		if getErr != nil {
			log.Fatalf("failed to create or fetch project: %v", getErr)
		}
		projectID = p.ID
	} else {
		fmt.Printf("Created project: %s (ID: %v)\n", project.Name, project.ID)
	}

	// 5. Add User to Project
	_, err = queries.AddUserToProject(ctx, db.AddUserToProjectParams{
		UserID:    userID,
		ProjectID: projectID,
		Role:      db.RoleOwner,
	})
	if err != nil {
		fmt.Printf("Note: User-Project link might already exist: %v\n", err)
	} else {
		fmt.Println("Linked user to project.")
	}

	fmt.Println("Database seeding completed successfully!")
}
