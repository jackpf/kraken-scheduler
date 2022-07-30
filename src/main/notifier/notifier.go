package notifier

type Notifier interface {
	Send(to string, subject string, body string) error
}
