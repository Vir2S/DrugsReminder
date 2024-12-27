# Drugs SMS Reminder REST API

## Prerequisites
Have installed:

Golang

Docker

PostgreSQL


## How to install
Run SQL query under postgresql shell:
```
CREATE DATABASE sms_reminder;

\c sms_reminder

CREATE TABLE patients (
id SERIAL PRIMARY KEY,
name VARCHAR(255) NOT NULL,
phone_number VARCHAR(20) NOT NULL,
email VARCHAR(255),
drug_name VARCHAR(255) NOT NULL,
dosage VARCHAR(255) NOT NULL,
reminder_time TIME NOT NULL,
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE INDEX idx_reminder_time ON patients(reminder_time);
```
Build and run docker container:

```
docker build -t sms-reminder .
docker run -d -p 8080:8080 --env-file .env sms-reminder
```

## How to use
To add new patient run from terminal:
```
curl -X POST http://localhost:8080/add-patient \
-H "Content-Type: application/json" \
-d '{
"name": "John Doe",
"phone_number": "1234567890",
"email": "john.doe@example.com",
"drug_name": "Vitamin C",
"dosage": "100 mg",
"reminder_time": "08:00:00"
}'
```

To send reminders:
```
curl http://localhost:8080/send-reminders
```
