package inituploadhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/initupload"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	tn        string
	expcod    int
	exbody    string
	data      *[]byte
	token     *string
	noAuthMid bool
}

const url = "https://localhost:8443"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
	statusU   = http.StatusUnauthorized
)

//nolint:lll
func getTestData() *[]testData {
	tok := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYzNTI1MTEsImlkIjoicXFxIn0.6O4koFaQNd4vdWWb2WsftHL0Ewb45_4-tng1qcPn980"

	tok1 := ""

	return &[]testData{
		{
			tn:     "1",
			expcod: statusU,
			exbody: "",
			data:   nil,
			token:  &tok,
		},
		{
			tn:     "2",
			expcod: statusU,
			exbody: "",
			data:   nil,
			token:  &tok1,
		},
		{
			tn:        "3",
			expcod:    statusBR,
			exbody:    "",
			data:      nil,
			token:     nil,
			noAuthMid: true,
		},
	}
}

func TestInitLoadHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(60 * time.Second)

	attr := &serverpa.ServerProcAttr{}

	err := attr.Init(true)
	if err != nil {
		t.Errorf("Init: %v", err)

		return
	}

	err = migrator.UseMigrations(attr)
	if err != nil {
		t.Errorf("Init: %v", err)

		return
	}

	path := "../../internal/client/auth/token.json"

	tok, err := authcfg.GetToken(path)
	if err != nil {
		t.Errorf("GetToken: %v", err)
	}

	Token := tok

	testCases := getTestData()

	initUploadH := initupload.NewHandler(attr.FIleService,
		attr.InitUploadAttr).InitUploadHandler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()
			req(t, &test, attr, initUploadH, Token)
		})
	}
}

func req(t *testing.T,
	test *testData,
	attr *serverpa.ServerProcAttr,
	handler func(
		writer http.ResponseWriter,
		req *http.Request,
	),
	token string,
) {
	t.Helper()

	tmp := make([]byte, 0)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url+"/api/user/initupload",
		bytes.NewReader(tmp))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	if test.token != nil {
		req.Header.Set("Authorization", *test.token)
	} else {
		req.Header.Set("Authorization", token)
	}

	newr := httptest.NewRecorder()
	router := mux.NewRouter()

	if !test.noAuthMid {
		router.Use(authmiddleware.AuthMiddleware(
			attr.AuthMidAttr))
	}

	router.HandleFunc("/api/user/initupload",
		handler)
	router.ServeHTTP(newr, req)
	status := newr.Code
	body, _ := io.ReadAll(newr.Body)

	assert.Equal(t,
		test.expcod,
		status, test.tn+": Response code didn't match expected")

	if test.exbody != "" {
		assert.JSONEq(t, test.exbody, string(body))
	}
}

//nolint:errchkjson
func GetIncorrectData() *[]byte {
	incd := &apim.IncorrectData{}
	incd.IncorrectData = "IncorrectData"
	marshal, _ := json.Marshal(incd)

	return &marshal
}
