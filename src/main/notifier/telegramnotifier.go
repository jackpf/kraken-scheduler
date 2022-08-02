package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type TelegramCredentials struct {
	Token  string
	ChatID int
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

	fmt.Printf("Sending %s to chat_id: %d", text, credentials.ChatID)

	var telegramApi string = "https://api.telegram.org/bot" + credentials.Token + "/sendMessage" // TODO get token from environment
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(credentials.ChatID)},
			"text":    {text},
		})

	if err != nil {
		fmt.Printf("An Error ocurred while posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		fmt.Printf("Error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)

	fmt.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

type TelegramNotifier struct {
	credentials TelegramCredentials
}

func NewTelegramNotifier(credentialsFile string) *TelegramNotifier {
	var credentials TelegramCredentials = readTelegramCredentials(credentialsFile)
	return &TelegramNotifier{credentials: credentials}
}

func (m TelegramNotifier) Send(subject string, body string) error {
	_, err := sendNotificationToTelegramChat(m.credentials, body)
	return err
}
