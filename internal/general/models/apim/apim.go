package apim

type CfgServer struct {
	ServerAddr string `json:"address"`
	DBDSN      string `json:"dbdsn"`
}

type InRegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InLoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
