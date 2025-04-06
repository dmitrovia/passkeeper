package euploaderattr

import "net/http"

type UploaderAttr struct {
	URL    string
	Client *http.Client
	Data   *[]byte
}

func (p *UploaderAttr) Init(
	eURL string,
	client *http.Client,
) {
	p.URL = eURL
	p.Client = client
}
