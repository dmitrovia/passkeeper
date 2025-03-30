package apim

type CfgServer struct {
	ServerAddr       string `json:"address"`
	DBDSN            string `json:"dbdsn"`
	FilesStoragePath string `json:"filesStoragePath"`
}

type CfgClient struct {
	ServerAddr           string `json:"address"`
	FilesSynchronizePath string `json:"filesSynchronizePath"`
	MetaPath             string `json:"metaPath"`
}

type InRegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InLoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
