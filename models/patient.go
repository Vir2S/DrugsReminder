package models

import "time"

type Patient struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	DrugName     string    `json:"drug_name"`
	Dosage       string    `json:"dosage"`
	ReminderTime time.Time `json:"reminder_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
