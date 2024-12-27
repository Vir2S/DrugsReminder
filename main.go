package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"DrugsReminder/config"
	"DrugsReminder/handlers"
	"DrugsReminder/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DBConnString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to establish database connection: %v", err)
	}

	log.Println("Successfully connected to the database")

	twilioClient := services.NewTwilioClient(
		cfg.TwilioConfig.AccountSID,
		cfg.TwilioConfig.AuthToken,
		cfg.TwilioConfig.FromNumber,
	)

	reminderHandler := handlers.NewReminderHandler(db, twilioClient)

	r := mux.NewRouter()
	r.HandleFunc("/add-patient", reminderHandler.AddPatient).Methods("POST")
	r.HandleFunc("/send-reminders", reminderHandler.SendReminders).Methods("GET")

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}
