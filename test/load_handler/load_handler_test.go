package loadhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/load"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	tn       string
	expcod   int
	exbody   string
	data     *[]byte
	token    *string
	Hash     *string
	FilePath *string
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

	fpath := "temp"

	thash1 := "2f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc51"

	return &[]testData{
		{
			tn:     "1",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,
			Hash:   nil,
		},
		{
			tn:     "2",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,
			Hash:   nil,
		},
		{
			tn:     "3",
			expcod: statusISE,
			exbody: "",
			data:   &tmp,
			token:  nil,

			Hash: nil,
		},
		{
			tn:     "4",
			expcod: statusBR,
			exbody: "",
			data:   incd,
			token:  nil,

			Hash: nil,
		},
		{
			tn:     "5",
			expcod: statusISE,
			exbody: "",
			data:   incd1,
			token:  nil,

			Hash: nil,
		},
		{
			tn:     "6",
			expcod: statusU,
			exbody: "",
			data:   nil,
			token:  &tok,

			Hash: nil,
		},
		{
			tn:     "7",
			expcod: statusU,
			exbody: "",
			data:   nil,
			token:  &tok1,

			Hash: nil,
		},
		{
			tn:     "8",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,

			Hash: nil,
		},
		{
			tn:     "9",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,

			Hash: &thash1,
		},
		{
			tn:     "10",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,

			Hash: nil,
		},
		{
			tn:     "11",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,

			Hash:     nil,
			FilePath: &fpath,
		},
		{
			tn:     "12",
			expcod: statusBR,
			exbody: "",
			data:   nil,
			token:  nil,

			Hash:     nil,
			FilePath: &fpath,
		},
	}
}

//nolint:funlen,cyclop
func TestLoadHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(60 * time.Second)

	attr := &serverpa.ServerProcAttr{}

	err := attr.Init()
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

	loadH := load.NewHandler(attr.MetaService,
		attr.LoadAttr).InitLoadHandler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()

			reqData, err := formReqBody(&test)
			if err != nil {
				t.Errorf("formReqBody: %v", err)

				return
			}

			var bodyReq []byte
			if test.data != nil {
				bodyReq = *test.data
			} else {
				bodyReq = *reqData
			}

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				url+"/api/user/load", bytes.NewReader(bodyReq))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			if test.token != nil {
				req.Header.Set("Authorization", *test.token)
			} else {
				req.Header.Set("Authorization", Token)
			}

			newr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Use(authmiddleware.AuthMiddleware(
				attr.AuthMidAttr))
			router.HandleFunc("/api/user/load",
				loadH)
			router.ServeHTTP(newr, req)
			status := newr.Code
			body, _ := io.ReadAll(newr.Body)

			assert.Equal(t,
				test.expcod,
				status, test.tn+": Response code didn't match expected")

			if test.exbody != "" {
				assert.JSONEq(t, test.exbody, string(body))
			}
		})
	}
}

func formReqBody(
	testd *testData,
) (*[]byte, error) {
	chunk := &chunckmeta.ChunkMeta{}

	chunk.Hash = testd.Hash
	chunk.FilePath = testd.FilePath

	marshal, err := json.Marshal(chunk)
	if err != nil {
		return nil, fmt.Errorf("formReqBody->Marshal: %w", err)
	}

	return &marshal, nil
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
