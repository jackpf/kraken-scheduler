package notifier

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func readCredentials(credentialsFile string) (*[]byte, error) {
	credentials, err := ioutil.ReadFile(credentialsFile)

	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	return &credentials, nil
}

// Request a Token from the web, then returns the retrieved Token.
func tokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-Token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code from the URL (code=..., paste the ... part): \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve Token from web: %v", err)
	}
	return tok
}

// Retrieves a Token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a Token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Debugf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth Token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Retrieve a Token, saves the Token
func getOrCreateToken(config *oauth2.Config, path string) *oauth2.Token {
	// The file Token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokenFile := fmt.Sprintf("%s/%s", path, "Token.json")
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		token = tokenFromWeb(config)
		saveToken(tokenFile, token)
	}

	return token
}

func MustNewGMailer(credentialsFile string, userId string) *GMailer {
	gmailer, err := NewGMailer(credentialsFile, userId)

	if err != nil {
		log.Fatal(err)
	}

	return gmailer
}

func NewGMailer(credentialsFile string, userId string) (*GMailer, error) {
	ctx := context.Background()

	credentials, err := readCredentials(credentialsFile)

	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved Token.json.
	config, err := google.ConfigFromJSON(*credentials, gmail.GmailSendScope, gmail.GmailComposeScope)

	if err != nil {
		return nil, err
	}

	credentialsFilePathAbsolute, err := filepath.Abs(credentialsFile)
	if err != nil {
		return nil, err
	}
	credentialsFilePathDir := filepath.Dir(credentialsFilePathAbsolute)

	token := getOrCreateToken(config, credentialsFilePathDir)
	client := config.Client(ctx, token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return nil, err
	}

	return &GMailer{service: service, userId: userId}, nil
}

type GMailer struct {
	service *gmail.Service
	userId  string
}

func (m GMailer) getUserEmailAddress() (*string, error) {
	profile, err := m.service.Users.GetProfile(m.userId).Do()

	if err != nil {
		return nil, err
	}

	return &profile.EmailAddress, nil
}

func (m GMailer) Send(subject string, body string) error {
	var message gmail.Message

	userEmailAddress, err := m.getUserEmailAddress()
	if err != nil {
		return err
	}

	emailTo := fmt.Sprintf("To: %s\r\n", *userEmailAddress)
	emailSubject := fmt.Sprintf("Subject: %s\r\n", subject)
	emailMime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	emailMessage := []byte(emailTo + emailSubject + emailMime + "\n" + body)

	message.Raw = base64.URLEncoding.EncodeToString(emailMessage)

	_, err = m.service.Users.Messages.Send(m.userId, &message).Do()

	return err
}
