package getsecretshandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecrets"
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
	tn        string
	expcod    int
	exbody    string
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
			token:  &tok,
		},
		{
			tn:     "2",
			expcod: statusU,
			exbody: "",
			token:  &tok1,
		},
		{
			tn:        "3",
			expcod:    statusBR,
			exbody:    "",
			token:     nil,
			noAuthMid: true,
		},
	}
}

func getTestData1() *[]testData {
	return &[]testData{
		{
			tn:     "4",
			expcod: statusISE,
			exbody: "",
			token:  nil,
		},
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

func TestGetSecretsHandler(t *testing.T) {
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
	testCases1 := getTestData1()

	getSecretsH := getsecrets.NewHandler(attr.SecretService,
		attr.GetSecretsAttr).GetSecretsHandler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()

			req(t, &test, attr, getSecretsH, Token)
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
			req(t, &test, attr, getSecretsH, Token)
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
		http.MethodGet,
		url+"/api/user/getsecrets", bytes.NewReader(tmp))
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

	router.HandleFunc("/api/user/getsecrets",
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

//nolint:errchkjson
func GetIncorrectDataWithCrypto(encKey *[]byte) *[]byte {
	incd := &apim.IncorrectData{}
	incd.IncorrectData = "IncorrectData"
	marshal, _ := json.Marshal(incd)

	encrypt, _ := rsa.Encrypt(&marshal, encKey)

	return encrypt
}
