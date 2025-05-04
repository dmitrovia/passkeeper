package einituploaderattr

import "net/http"

type InitUploadAttr struct {
	URL    string
	Client *http.Client
	Token  string
}

func (p *InitUploadAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
