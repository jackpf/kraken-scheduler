package notifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TelegramCredentials struct {
	token  string
	chatId int
}

func readTelegramCredentials(credentialsFile string) TelegramCredentials {
	// Let's first read the `telegram-credentials.json` file
	content, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var credentials TelegramCredentials
	err = json.Unmarshal(content, &credentials)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return credentials
}

// sendNotificationToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendNotificationToTelegramChat(credentials TelegramCredentials, text string) (string, error) {

	log.Debug("Sending %s to chat_id: %d", text, credentials.chatId)

	var telegramApi string = "https://api.telegram.org/bot" + credentials.token + "/sendMessage" // TODO get token from environment
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(credentials.chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Debug("An Error ocurred while posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Debug("Error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)

	log.Debug("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
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
