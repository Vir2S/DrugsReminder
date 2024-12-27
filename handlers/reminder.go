package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"DrugsReminder/models"
	"DrugsReminder/services"
)

type ReminderHandler struct {
	DB           *sql.DB
	TwilioClient *services.TwilioClient
}

func NewReminderHandler(db *sql.DB, twilioClient *services.TwilioClient) *ReminderHandler {
	return &ReminderHandler{
		DB:           db,
		TwilioClient: twilioClient,
	}
}

func (h *ReminderHandler) SendReminders(w http.ResponseWriter, r *http.Request) {
	// Поточний час без дати
	currentTime := time.Now().Format("15:04:05")

	query := `SELECT id, name, phone_number, drug_name, dosage
              FROM patients
              WHERE reminder_time = $1`

	rows, err := h.DB.Query(query, currentTime)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.PhoneNumber, &p.DrugName, &p.Dosage); err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		message := fmt.Sprintf(
			"Hello %s! This is a reminder to take your medicine: %s. Dosage: %s.",
			p.Name, p.DrugName, p.Dosage)

		sendReq := services.SendSMSRequest{
			To:   p.PhoneNumber,
			Body: message,
		}

		sendResp, err := h.TwilioClient.SendSMS(sendReq)
		if err != nil {
			log.Printf("Failed to send SMS to %s: %v", p.PhoneNumber, err)
			continue
		}

		log.Printf("Reminder sent to patient %s, Message SID: %s", p.PhoneNumber, sendResp.MessageID)
	}

	// Перевіряємо, чи виникли помилки при ітерації по рядках
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Відповідаємо, що все добре
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reminders sent"))
}

// (Опціонально) можете зробити окрему функцію для побудови повідомлення
func buildReminderMessage(p models.Patient) string {
	return fmt.Sprintf("Hello %s! This is a reminder to take your medicine: %s. Dosage: %s.",
		p.Name, p.DrugName, p.Dosage)
}
