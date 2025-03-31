package entity

type EmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewEmailSender(host, username, password, from string, port int) *EmailSender {
	return &EmailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}
