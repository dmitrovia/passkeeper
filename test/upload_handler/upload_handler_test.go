package uploadhandler_test

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
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/aes256"
	"github.com/dmitrovia/passkeeper/internal/general/compress"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/upload"
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
	expcod       int
	exbody       string
	data         *[]byte
	token        *string
	FileName     string
	OrigFileName string
	Index        int
	Hash         *string
	DataMeta     *[]byte
	FilePath     *string
	noAuthMid    bool
}

const url = "https://localhost:8443"

const path = "../../internal/client/crypto/keys/public.pem"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
	statusU   = http.StatusUnauthorized
)

const strLen = 5

//nolint:lll,funlen
func getTestData(encKey *[]byte) *[]testData {
	tmp := make([]byte, 0)

	incd := GetIncorrectData()

	incd1 := GetIncorrectDataWithCrypto(encKey)

	tok := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYzNTI1MTEsImlkIjoicXFxIn0.6O4koFaQNd4vdWWb2WsftHL0Ewb45_4-tng1qcPn980"

	tok1 := ""

	dmeta := []byte(randomString())

	fpath := "temp"

	thash := "2f282b84e7e608d5852449ed940bfc51"

	thash1 := "2f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc512f282b84e7e608d5852449ed940bfc51"

	return &[]testData{
		{
			tn:           "1",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			FileName:     "",
			OrigFileName: "",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &tmp,
		},
		{
			tn:           "2",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			FileName:     "",
			OrigFileName: "",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "3",
			expcod:       statusISE,
			exbody:       "",
			data:         &tmp,
			token:        nil,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "4",
			expcod:       statusBR,
			exbody:       "",
			data:         incd,
			token:        nil,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "5",
			expcod:       statusISE,
			exbody:       "",
			data:         incd1,
			token:        nil,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "6",
			expcod:       statusU,
			exbody:       "",
			data:         nil,
			token:        &tok,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "7",
			expcod:       statusU,
			exbody:       "",
			data:         nil,
			token:        &tok1,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "8",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:           "9",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			FileName:     "upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash1,
			DataMeta:     &dmeta,
		},
		{
			tn:           "10",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			FileName:     "upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0upload_test1.chunk.0",
			OrigFileName: "upload_test1",
			Index:        99,
			Hash:         &thash,
			DataMeta:     &dmeta,
		},
		{
			tn:       "11",
			expcod:   statusBR,
			exbody:   "",
			data:     nil,
			token:    nil,
			FileName: "upload_test1",
			Index:    99,
			Hash:     &thash,
			FilePath: &fpath,
			DataMeta: &dmeta,
		},
		{
			tn:       "12",
			expcod:   statusBR,
			exbody:   "",
			data:     nil,
			token:    nil,
			FileName: "upload_test1",
			Index:    99,
			Hash:     nil,
			FilePath: &fpath,
			DataMeta: &dmeta,
		},
		{
			tn:        "13",
			expcod:    statusBR,
			exbody:    "",
			data:      nil,
			token:     nil,
			FileName:  "upload_test1",
			Index:     99,
			Hash:      nil,
			FilePath:  &fpath,
			DataMeta:  &dmeta,
			noAuthMid: true,
		},
	}
}

func getTestData1() *[]testData {
	tmp := "temp"
	dmeta := []byte(randomString())

	return &[]testData{
		{
			tn:           "14",
			expcod:       statusISE,
			exbody:       "",
			data:         nil,
			token:        nil,
			FileName:     "upload_test1",
			Index:        99,
			Hash:         &tmp,
			FilePath:     nil,
			DataMeta:     &dmeta,
			OrigFileName: "upload_test1",
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

//nolint:funlen
func TestUploadHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(60 * time.Second)

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

	uploadH := upload.NewHandler(attr.FIleService,
		attr.MetaService,
		attr.UploadAttr).UploadHandler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()
			req(t, &test, uploadH, Token, attr)
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
			req(t, &test, uploadH, Token, attr)
		})
	}
}

//nolint:funlen
func req(t *testing.T,
	test *testData,
	handler func(
		writer http.ResponseWriter,
		req *http.Request,
	),
	token string,
	attr *serverpa.ServerProcAttr,
) {
	t.Helper()

	reqData, err := formReqBody(test)
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
		url+"/api/user/upload", bytes.NewReader(bodyReq))
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

	router.HandleFunc("/api/user/upload",
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
) (*[]byte, error) {
	chunk := &chunckmeta.ChunkMeta{}

	chunk.FileName = &testd.FileName
	chunk.OrigFileName = &testd.OrigFileName
	chunk.Hash = testd.Hash
	chunk.Index = &testd.Index
	chunk.Data = testd.DataMeta
	chunk.FilePath = testd.FilePath

	err := compressAndEncrypt(chunk, chunk.Data)
	if err != nil {
		return nil, fmt.Errorf("toChunk->compressAndEncrypt: %w",
			err)
	}

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

func randomString() string {
	letters := []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, strLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func compressAndEncrypt(
	chunk *chunckmeta.ChunkMeta,
	chBytes *[]byte,
) error {
	Aes256key := "gLn2xbvYdE556NYzlaQKPwIFuSAY6NXJ"
	Aes256keyBytes := []byte(Aes256key)

	compressData, err := compress.DeflateCompress(
		*chBytes)
	if err != nil {
		return fmt.Errorf("toChunk->DC: %w", err)
	}

	dec, err := aes256.Encrypt(&compressData,
		&Aes256keyBytes)
	if err != nil {
		return fmt.Errorf("PRASF->aes256Decrypt: %w", err)
	}

	chunk.Data = dec

	return nil
}
