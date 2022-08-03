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
	token  string
	chatId int
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

func apiCall(path string, parameters url.Values) (*string, error) {
	response, err := http.PostForm(path, parameters)

	if err != nil {
		log.Debug("An Error occurred while posting text to the chat: %s", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Debug("Error in parsing telegram answer %s", errRead.Error())
		return nil, err
	}
	bodyString := string(bodyBytes)

	log.Debug("Body of Telegram Response: %s", bodyString)

	return &bodyString, nil
}

// sendNotificationToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendNotificationToTelegramChat(credentials TelegramCredentials, text string) (*string, error) {
	log.Debugf("Sending %s to chat_id: %d", text, credentials.chatId)

	return apiCall(fmt.Sprintf(telegramApiSendMessagePath, credentials.token), url.Values{
		"chat_id": {strconv.Itoa(credentials.chatId)},
		"text":    {text},
	})
}

type TelegramNotifier struct {
	credentials TelegramCredentials
}

func NewTelegramNotifier(credentialsFile string) *TelegramNotifier {
	credentials := readTelegramCredentials(credentialsFile)
	return &TelegramNotifier{credentials: credentials}
}

func (m TelegramNotifier) Send(subject string, body string) error {
	_, err := sendNotificationToTelegramChat(m.credentials, body)
	return err
}
