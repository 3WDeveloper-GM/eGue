package indexer

type Email struct {
	MessageID              string `json:"message-id"`
	Date                   string `json:"date"`
	From                   string `json:"from"`
	To                     string `json:"to"`
	Subject                string `json:"subject"`
	MimeVersion            string `json:"mime-version"`
	ContentType            string `json:"content-type"`
	ContentTranferEncoding string `json:"content-transfer-encoding"`
	XFrom                  string `json:"x-from"`
	Xto                    string `json:"x-to"`
	Xcc                    string `json:"x-cc"`
	Xbcc                   string `json:"x-bcc"`
	XFolder                string `json:"x-folder"`
	XOrigin                string `json:"x-origin"`
	XFileName              string `json:"x-filename"`
	CC                     string `json:"cc"`
	BCC                    string `json:"bcc"`
	Bod                    string `json:"body"`
}

type Index struct {
}
