package notifications

type Notification interface {
	Subject() string
	Body() string
}
