package uploaderattrs

import "net/http"

type UploaderAttr struct {
	ServerURL string
	Client    *http.Client
	Data      *[]byte
}

func (p *UploaderAttr) Init(
	serverURL string,
	client *http.Client,
) {
	p.ServerURL = serverURL
	p.Client = client
}
