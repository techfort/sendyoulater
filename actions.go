package sendyoulater

type Action struct{}

type EmailAction struct {
	Action
	UserID  string
	To      string
	Subject string
	Body    string
}
