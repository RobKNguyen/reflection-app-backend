package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "reflection-app/internal/database"
    "reflection-app/internal/handlers"
    "reflection-app/internal/repository"
    "reflection-app/internal/service"
    "reflection-app/internal/test"
)

func main() {
    // Load environment variables from .env file (for local development)
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    
    // Initialize database
    db, err := database.NewConnection()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Create tables
    if err := database.CreateTables(db); err != nil {
        log.Fatal("Failed to create tables:", err)
    }

    // Check if reflection_tracking table exists
    var tableExists bool
    err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'reflection_tracking')").Scan(&tableExists)
    if err != nil {
        log.Printf("Error checking reflection_tracking table: %v", err)
    } else if !tableExists {
        log.Println("reflection_tracking table does not exist, creating it...")
        trackingTable := `
        CREATE TABLE IF NOT EXISTS reflection_tracking (
            id SERIAL PRIMARY KEY,
            reflection_id INTEGER NOT NULL REFERENCES reflections(id) ON DELETE CASCADE,
            user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            reflected_date DATE NOT NULL DEFAULT CURRENT_DATE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(reflection_id, user_id, reflected_date)
        )`
        
        if _, err := db.Exec(trackingTable); err != nil {
            log.Printf("Error creating reflection_tracking table: %v", err)
        } else {
            log.Println("Created reflection_tracking table")
        }
    } else {
        log.Println("reflection_tracking table exists")
    }
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    reflectionRepo := repository.NewReflectionRepository(db)
    actionRepo := repository.NewActionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)  // Add this
	friendshipRepo := repository.NewFriendshipRepository(db)
	reactionRepo := repository.NewReactionRepository(db)
    
    // Initialize services
	authService := service.NewAuthService(userRepo)  // Add this
    userService := service.NewUserService(userRepo)
    reflectionService := service.NewReflectionService(reflectionRepo)
    actionService := service.NewActionService(actionRepo)
	categoryService := service.NewCategoryService(categoryRepo)  // Add this
	friendshipService := service.NewFriendshipService(friendshipRepo)
	reactionService := service.NewReactionService(reactionRepo)
    
    // Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)  // Add this
    userHandler := handlers.NewUserHandler(userService)
    reflectionHandler := handlers.NewReflectionHandler(reflectionService)
    actionHandler := handlers.NewActionHandler(actionService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)  // Add this
	friendshipHandler := handlers.NewFriendshipHandler(friendshipService)
	reactionHandler := handlers.NewReactionHandler(reactionService)
    testHandler := test.NewTestHandler()
    
    // Setup routes
    r := mux.NewRouter()
    
    // Enable CORS for frontend development
    r.Use(corsMiddleware)
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()

	// Auth routes - ADD THESE
    api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
    api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
    api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
    api.HandleFunc("/auth/me", authHandler.GetCurrentUser).Methods("GET")
    
    // User routes
    api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
    api.HandleFunc("/users/search", userHandler.SearchUsers).Methods("GET")  // Move this BEFORE the {id} route
    api.HandleFunc("/users/username/{username}", userHandler.GetUserByUsername).Methods("GET")
    api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
    api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
    api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
    
    // Reflection routes
    api.HandleFunc("/reflections", reflectionHandler.CreateReflection).Methods("POST")
    api.HandleFunc("/reflections", reflectionHandler.GetReflections).Methods("GET")
    api.HandleFunc("/reflections/user/{userId}", reflectionHandler.GetUserReflections).Methods("GET")
    api.HandleFunc("/reflections/{id}", reflectionHandler.GetReflection).Methods("GET")
    api.HandleFunc("/reflections/{id}", reflectionHandler.UpdateReflection).Methods("PUT")
    api.HandleFunc("/reflections/{id}", reflectionHandler.DeleteReflection).Methods("DELETE")
    api.HandleFunc("/reflections/{id}/reflect", reflectionHandler.TrackReflection).Methods("POST")
    
    // Analytics routes
    api.HandleFunc("/analytics/reflection-tracking", reflectionHandler.GetReflectionTrackingAnalytics).Methods("GET")
    api.HandleFunc("/analytics/reflection-tracking-by-category", reflectionHandler.GetReflectionTrackingByCategory).Methods("GET")
    
    // Action routes
    api.HandleFunc("/actions", actionHandler.CreateAction).Methods("POST")
    api.HandleFunc("/reflections/{reflectionId}/actions", actionHandler.GetActionsByReflection).Methods("GET")
    api.HandleFunc("/actions/{id}", actionHandler.GetAction).Methods("GET")
    api.HandleFunc("/actions/{id}", actionHandler.UpdateAction).Methods("PUT")
    api.HandleFunc("/actions/{id}/status", actionHandler.UpdateActionStatus).Methods("PATCH", "PUT")
    api.HandleFunc("/actions/{id}/complete", actionHandler.CompleteAction).Methods("PATCH")
    api.HandleFunc("/actions/{id}", actionHandler.DeleteAction).Methods("DELETE")

	// Category routes - Add these
    api.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET")
    api.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
    api.HandleFunc("/categories/{categoryId}/subcategories", categoryHandler.GetSubCategories).Methods("GET")
    api.HandleFunc("/categories/{categoryId}/subcategories", categoryHandler.CreateSubCategory).Methods("POST")
    api.HandleFunc("/subcategories/{id}", categoryHandler.UpdateSubCategory).Methods("PUT")
    api.HandleFunc("/subcategories/{id}", categoryHandler.DeleteSubCategory).Methods("DELETE")
    api.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
    api.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")
    
    // Social feed routes
    api.HandleFunc("/feed/friends", reflectionHandler.GetFriendsFeed).Methods("GET")
    api.HandleFunc("/feed/friend/{friendId}", reflectionHandler.GetFriendReflections).Methods("GET")
    api.HandleFunc("/feed/friend-by-username", reflectionHandler.GetFriendReflectionsByUsername).Methods("GET")
    
    // Friendship routes
    api.HandleFunc("/friends/request", friendshipHandler.SendFriendRequest).Methods("POST")
    api.HandleFunc("/friends/accept", friendshipHandler.AcceptFriendRequest).Methods("POST")
    api.HandleFunc("/friends/reject", friendshipHandler.RejectFriendRequest).Methods("POST")
    api.HandleFunc("/friends", friendshipHandler.GetFriendsList).Methods("GET")
    api.HandleFunc("/friends/pending", friendshipHandler.GetPendingRequests).Methods("GET")
    api.HandleFunc("/friends", friendshipHandler.RemoveFriend).Methods("DELETE")
    
    // Reaction routes
    api.HandleFunc("/reactions", reactionHandler.AddReaction).Methods("POST")
    api.HandleFunc("/reactions", reactionHandler.RemoveReaction).Methods("DELETE")
    api.HandleFunc("/reactions", reactionHandler.GetReactionsForReflection).Methods("GET")
    api.HandleFunc("/reactions/counts", reactionHandler.GetReactionCountsForReflection).Methods("GET")
    api.HandleFunc("/reactions/user", reactionHandler.GetUserReactionForReflection).Methods("GET")
    api.HandleFunc("/reactions/prompts", reactionHandler.GetReactionPrompts).Methods("GET")

    // Test routes
    api.HandleFunc("/test/sql", testHandler.SQLTest).Methods("GET", "POST")
    api.HandleFunc("/test/echo", testHandler.Echo).Methods("GET", "POST")
    api.HandleFunc("/test/create", testHandler.Create).Methods("POST")
    
    // Health check endpoint
    api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")
    
    // Test endpoint for debugging
    api.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Test endpoint called - Method: %s, URL: %s", r.Method, r.URL.String())
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Backend is working!", "method": "` + r.Method + `", "url": "` + r.URL.String() + `"}`))
    }).Methods("GET")
    
    // Serve static files (for frontend)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
    
    // Get port from environment variable (for Render deployment)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

// CORS middleware for development
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := os.Getenv("CORS_ORIGIN")
        if origin == "" {
            origin = "http://localhost:3000"
        }
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}