package notificationtemplates

type NotificationTemplate interface {
	Subject() string
	Body() string
}
