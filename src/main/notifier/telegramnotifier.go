package notifier

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// sendNotificationToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendNotificationToTelegramChat(chatId int, text string) (string, error) {

	fmt.Printf("Sending %s to chat_id: %d", text, chatId)

	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage" // TODO get token from environment
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
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
