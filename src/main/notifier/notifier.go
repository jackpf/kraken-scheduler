package notifier

type Notifier interface {
	Send(subject string, body string) error
}
