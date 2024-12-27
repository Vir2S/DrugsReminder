package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseUrl = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"
)

type TwilioClient struct {
	accountSID string
	authToken  string
	fromNumber string
}

func NewTwilioClient(accountSID, authToken, fromNumber string) *TwilioClient {
	return &TwilioClient{
		accountSID: accountSID,
		authToken:  authToken,
		fromNumber: fromNumber,
	}
}

type SendSMSRequest struct {
	To   string `json:"To"`
	From string `json:"From"`
	Body string `json:"Body"`
}

type SendSMSResponse struct {
	MessageID string `json:"message_id"`
}

func (c *TwilioClient) SendSMS(req SendSMSRequest) (SendSMSResponse, error) {
	url := fmt.Sprintf(baseUrl, c.accountSID)

	// Формуємо Basic Auth (логін:пароль)
	auth := fmt.Sprintf("%s:%s", c.accountSID, c.authToken)
	// Тут краще за все використати Base64
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	// Серіалізуємо тіло запиту в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("request marshaling error: %v", err)
	}

	// Створюємо http-запит (використовуємо іншу змінну замість `req`)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Встановлюємо заголовки
	httpReq.Header.Set("Authorization", "Basic "+encodedAuth)
	httpReq.Header.Set("Content-Type", "application/json")

	// Відправляємо запит
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Читаємо відповідь
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("error reading response: %v", err)
	}

	// Перевіряємо статус
	if resp.StatusCode != http.StatusOK {
		return SendSMSResponse{}, fmt.Errorf("Twilio API error: %s", string(body))
	}

	// Парсимо JSON-відповідь Twilio
	var sendResp SendSMSResponse
	err = json.Unmarshal(body, &sendResp)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("parsing error: %v", err)
	}

	return sendResp, nil
}
