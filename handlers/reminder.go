package handlers

import (
	"database/sql"
	"encoding/json"
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

func (h *ReminderHandler) AddPatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if patient.Name == "" || patient.PhoneNumber == "" || patient.DrugName == "" || patient.Dosage == "" || patient.ReminderTime.IsZero() {
		http.Error(w, "Не все обязательные поля заполнены", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO patients (name, phone_number, email, drug_name, dosage, reminder_time, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`

	err = h.DB.QueryRow(query, patient.Name, patient.PhoneNumber, patient.Email, patient.DrugName, patient.Dosage, patient.ReminderTime.Format("15:04:05")).Scan(&patient.ID)
	if err != nil {
		log.Printf("Ошибка вставки в базу данных: %v", err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

func (h *ReminderHandler) SendReminders(w http.ResponseWriter, r *http.Request) {
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

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reminders sent"))
}

func buildReminderMessage(p models.Patient) string {
	return fmt.Sprintf("Hello %s! This is a reminder to take your medicine: %s. Dosage: %s.",
		p.Name, p.DrugName, p.Dosage)
}
