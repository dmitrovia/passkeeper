package einitsingleloadattr

import "net/http"

type InitSingleLoadAttr struct {
	URL    string
	Client *http.Client
	Token  string
	Data   *[]byte
}

func (p *InitSingleLoadAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
