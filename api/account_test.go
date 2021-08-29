package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/bank-demo/db/mock"
	db "github.com/bank-demo/db/sqlc"
	"github.com/bank-demo/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// ----- after mock db
// write unit test for 
// GetAccount API
// CreateAccount API
// ListAccount API
func TestGetAccountAPI(t *testing.T) {
	// need to create an account
	account := randomAccount()
	// need to create a new mock store
	//using  mockdb.NewMockStore() generated function.
	//It expects a gomock.Controller object as input,
	//so we have to create this controller by calling gomock.
	//NewController and pass in the testing.T object.

	// final: table-driven test
	// cases : 
	// get account success
	// account not found in db - no row in db
	// internal error - error occurred on connection between server and db.
	// invalid ID - bad request, invalid ID doesn't satify binding condition
	testCases := []struct{
		name string
		accountID int64
		// for different scenario to be built
		buildStub func(store *mockdb.MockStore)
		// check the output API
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// get account success
		{
			name: "OK",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// Account Not found
		{
			name: "NotFound",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Accounts{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		// Internal Error
		{
			name: "InternalError",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Accounts{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// Invalid ID
		{
			name: "InvalidID",
			accountID: 0,
			buildStub: func(store *mockdb.MockStore) {
	// because the invalid ID, the GetAccount should not call the handler
	store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// remember leetcode test case
	// use loop to iterate each test case
	for i := range testCases{
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	//create a new store by calling mockdb.NewMockStore() with this input controller.
	// store := mockdb.NewMockStore(ctrl)
	// build stub for mock store
	// specify what values of these 2 parameters we expect this function to be called with.
	// The first context argument could be any value,
	// The second argument should equal to the ID of the random account has created above.

	// store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID))
	// Now this stub definition can be translated as:
	// expect the GetAccount() function of the store to be called with any context
	// and this specific account ID arguments.
	//specify how many times this function should be called using the Times() function
	// use the Return() function to tell gomock to return some specific values
	// whenever the GetAccount() function is called.
	// in this case, we want it to return the account object and a nil error.
	////store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	// ---- stub for mock Store is built, Start to test Sever
	// server := NewServer(store)
	// use the recording feature of httptest package to record the response
	// of the API request. think of recorder is a fake respoonse
	// recorder := httptest.NewRecorder()
	// declare the url path of the API we want to call
	// url := fmt.Sprintf("/accounts/%d", account.ID)
	// create a new HTTP Request with method GET to that URL, and return request object or err
	//	request, err := http.NewRequest(http.MethodGet, url, nil)
	// test for err and request
	// require.NoError(t, err)
	// if no error, test the API Server
	// use the recorder and request to call server.router.ServeHTTP()
	// server.router.ServeHTTP(recorder, request)
	// require.Equal(t, http.StatusOK, recorder.Code)
	// // check the body of response
	// // need to create another function to get the body of response
	// // the body of response is stored in recorder
	// requireBodyMatchAccount(t, recorder.Body, account)
}

// test Create Account API

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 10000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
// body is the body of response, account is the object to be compared
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// gotAccount store the account object we got from body data
	var gotAccount db.Accounts
	// unmarshall data to gotAccount object
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	// gotAccount must be equal to input account
	require.Equal(t, account, gotAccount)
}
