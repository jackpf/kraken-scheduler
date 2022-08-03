package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	telegramApiBasePath        = "https://api.telegram.org/bot"
	telegramApiSendMessagePath = telegramApiBasePath + "%s/sendMessage"
)

type TelegramCredentials struct {
	Token  string `json:"token"`
	ChatId int    `json:"chatId"`
}

func readTelegramCredentials(credentialsFile string) TelegramCredentials {
	content, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var credentials TelegramCredentials
	if err = json.Unmarshal(content, &credentials); err != nil {
		log.Fatal("Unable to parse credentials file: ", err)
	}

	return credentials
}

func apiCall(path string, parameters url.Values) (map[string]interface{}, error) {
	fmt.Printf("Path: %s", path)
	response, err := http.PostForm(path, parameters)

	if err != nil {
		log.Printf("An Error occurred while posting text to the chat: %s", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Debugf("Error in parsing telegram answer %s", errRead.Error())
		return nil, err
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &jsonBody)
	if err != nil {
		return nil, err
	}

	if ok, hasOk := jsonBody["ok"]; !hasOk || ok.(bool) == false {
		return nil, fmt.Errorf("telegram API call failed: %+v", jsonBody)
	}

	log.Printf("Telegram Response: %+v", jsonBody)

	return jsonBody, nil
}

// sendNotificationToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendNotificationToTelegramChat(credentials TelegramCredentials, text string) (map[string]interface{}, error) {
	log.Debugf("Sending %s to chat_id: %d", text, credentials.ChatId)

	return apiCall(fmt.Sprintf(telegramApiSendMessagePath, credentials.Token), url.Values{
		"chat_id": {strconv.Itoa(credentials.ChatId)},
		"text":    {text},
	})
}

type TelegramNotifier struct {
	credentials TelegramCredentials
}

func MustNewTelegramNotifier(credentialsFile string) *TelegramNotifier {
	notifier, err := NewTelegramNotifier(credentialsFile)

	if err != nil {
		log.Fatal(err.Error())
	}

	return notifier
}

func NewTelegramNotifier(credentialsFile string) (*TelegramNotifier, error) {
	credentials := readTelegramCredentials(credentialsFile)

	if credentials.Token == "" {
		return nil, fmt.Errorf("telegram token must be configured")
	}
	if credentials.ChatId == 0 {
		return nil, fmt.Errorf("telegram chatId must be configured")
	}

	return &TelegramNotifier{credentials: credentials}, nil
}

func (m TelegramNotifier) Send(subject string, body string) error {
	_, err := sendNotificationToTelegramChat(m.credentials, fmt.Sprintf("%s\n\n%s", subject, body))
	return err
}
