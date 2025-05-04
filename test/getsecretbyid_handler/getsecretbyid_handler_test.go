package getsecretbyidhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecretbyid"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware/authmiddlewareattr"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
	"github.com/dmitrovia/passkeeper/internal/server/service/authservice"
	"github.com/dmitrovia/passkeeper/internal/server/storage/userstorage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	tn           string
	inIdentifier string
	expcod       int
	exbody       string
	data         *[]byte
	token        *string
	noAuthMid    bool
	noEncr       bool
}

const url = "https://localhost:8443"

const path = "../../internal/client/crypto/keys/public.pem"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
	statusU   = http.StatusUnauthorized
)

//nolint:lll,funlen
func getTestData(encKey *[]byte) *[]testData {
	tmp := make([]byte, 0)

	incd := GetIncorrectData()

	incd1 := GetIncorrectDataWithCrypto(encKey)

	tok := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYzNTI1MTEsImlkIjoicXFxIn0.6O4koFaQNd4vdWWb2WsftHL0Ewb45_4-tng1qcPn980"

	tok1 := ""

	return &[]testData{
		{
			tn:           "1",
			inIdentifier: "test" + randomString(),
			expcod:       statusISE,
			exbody:       "",
			data:         &tmp,
			token:        nil,
		},
		{
			tn:           "2",
			inIdentifier: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest" + randomString(),
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
		},
		{
			tn:           "3",
			inIdentifier: "test" + randomString(),
			expcod:       statusISE,
			exbody:       "",
			data:         incd,
			token:        nil,
		},
		{
			tn:           "4",
			inIdentifier: "test" + randomString(),
			expcod:       statusBR,
			exbody:       "",
			data:         incd1,
			token:        nil,
		},
		{
			tn:           "5",
			inIdentifier: "",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
		},
		{
			tn:           "6",
			inIdentifier: "",
			expcod:       statusU,
			exbody:       "",
			data:         nil,
			token:        &tok,
		},
		{
			tn:           "7",
			inIdentifier: "",
			expcod:       statusU,
			exbody:       "",
			data:         nil,
			token:        &tok1,
		},
		{
			tn:           "8",
			inIdentifier: "test",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			noAuthMid:    true,
		},
		{
			tn:           "9",
			inIdentifier: "test",
			expcod:       statusISE,
			exbody:       "",
			data:         nil,
			token:        nil,
			noEncr:       true,
		},
	}
}

func getTestData1() *[]testData {
	return &[]testData{
		{
			tn:           "11",
			inIdentifier: "test",
			expcod:       statusISE,
			exbody:       "",
			data:         nil,
			token:        nil,
		},
	}
}

//nolint:funlen
func TestGetSByIdHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	// time.Sleep(60 * time.Second)

	attr := &serverpa.ServerProcAttr{}

	err := attr.Init(true)
	if err != nil {
		t.Errorf("Init: %v", err)

		return
	}

	encKey, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("ReadFile: %v", err)

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

	testCases := getTestData(&encKey)
	testCases1 := getTestData1()

	getSecretByIDH := getsecretbyid.NewHandler(
		attr.SecretService,
		attr.GetSecretByIDAttr).GetSecretByIDHadnler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()
			req(t, &test, attr, getSecretByIDH, Token, encKey)
		})
	}

	attr.PgxConn.Close()

	err = newConn(attr)
	if err != nil {
		t.Errorf("newConn: %v", err)
	}

	for _, test := range *testCases1 {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()
			req(t, &test, attr, getSecretByIDH, Token, encKey)
		})
	}
}

func newConn(attr *serverpa.ServerProcAttr) error {
	ctxDB, cancel := context.WithTimeout(
		context.Background(), attr.Dbtimeout)
	defer cancel()

	dbConn, err := pgxpool.New(ctxDB,
		attr.DBDSN)
	if err != nil {
		return fmt.Errorf("SetPgxPool->pgxpool.New: %w", err)
	}

	UserStorage := &userstorage.UserStorage{}
	UserStorage.Initiate(dbConn)

	attr.AuthService = authservice.NewAuthService(
		UserStorage)

	attr.AuthMidAttr = &authmiddlewareattr.AuthMiddlewareAttr{}
	attr.AuthMidAttr.Init(attr.ZapLogger,
		attr.AuthService, attr.Dbtimeout, attr.SecretAuth)

	return nil
}

//nolint:funlen
func req(t *testing.T,
	test *testData,
	attr *serverpa.ServerProcAttr,
	handler func(
		writer http.ResponseWriter,
		req *http.Request,
	),
	token string,
	encKey []byte,
) {
	t.Helper()

	reqData, err := formReqBody(test, &encKey)
	if err != nil {
		t.Errorf("formReqBody: %v", err)

		return
	}

	var bodyReq *bytes.Reader
	if test.data != nil {
		bodyReq = bytes.NewReader(*test.data)
	} else {
		bodyReq = bytes.NewReader(*reqData)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		url+"/api/user/getsecretbyid", bodyReq)
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

	router.HandleFunc("/api/user/getsecretbyid",
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

func formReqBody(
	testd *testData,
	encKey *[]byte,
) (*[]byte, error) {
	outAttr := &apim.InGetSecretByID{}
	outAttr.Identifier = testd.inIdentifier

	marshal, err := json.Marshal(outAttr)
	if err != nil {
		return nil, fmt.Errorf("formReqBody->Marshal: %w", err)
	}

	if testd.noEncr {
		return &marshal, nil
	}

	encrypt, err := rsa.Encrypt(&marshal, encKey)
	if err != nil {
		return nil, fmt.Errorf("formReqBody->Encrypt: %w", err)
	}

	return encrypt, nil
}

//nolint:errchkjson
func GetIncorrectData() *[]byte {
	incd := &apim.IncorrectData{}
	incd.IncorrectData = "IncorrectData"
	marshal, _ := json.Marshal(incd)

	return &marshal
}

//nolint:errchkjson
func GetIncorrectDataWithCrypto(encKey *[]byte) *[]byte {
	incd := &apim.IncorrectData{}
	incd.IncorrectData = "IncorrectData"
	marshal, _ := json.Marshal(incd)

	encrypt, _ := rsa.Encrypt(&marshal, encKey)

	return encrypt
}

func randomString() string {
	letters := []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
