package domain

const (
	MailDir = "/enron_mail_20110402/maildir"
	Index   = "mails"
)

type Email struct {
	Date    string `json:"date"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	CC      string `json:"cc"`
	BCC     string `json:"bcc"`
	Body    string `json:"body"`
}
