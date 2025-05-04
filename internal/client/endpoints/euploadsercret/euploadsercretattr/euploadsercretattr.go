package euploadsercretattr

import "net/http"

type UploadSecretAttr struct {
	URL    string
	Client *http.Client
	Token  string
	Data   *[]byte
}

func (p *UploadSecretAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
