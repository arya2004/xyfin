package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/arya2004/xyfin/db/mock"
	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T){

	testCases := []struct {
		name         string
		accountID    int64
		buildStubs   func(store *mockdb.MockStore)
		expectStatus int
	}{
		{
			name:      "OK",
			accountID: 1,
			buildStubs: func(store *mockdb.MockStore) {
				const id int64 = 1
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(id)).
					Return(randomAccount(id), nil).
					Times(1)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:      "NotFound",
			accountID: 2,
			buildStubs: func(store *mockdb.MockStore) {
				const id int64 = 2
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(id)).
					Return(db.Account{}, sql.ErrNoRows).
					Times(1)
			},
			expectStatus: http.StatusNotFound,
		},
		{
			name:      "InternalError",
			accountID: 3,
			buildStubs: func(store *mockdb.MockStore) {
				const id int64 = 3
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(id)).
					Return(db.Account{}, sql.ErrConnDone).
					Times(1)
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectStatus: http.StatusBadRequest,
		},
	}


	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, tc.expectStatus, recorder.Code)
		})
	}
}


func randomAccount(id int64) db.Account {
	return db.Account{
		ID:        id,
		Owner:     utils.RandomOwner(),
		Balance:   utils.RandomMoney(),
		Currency:  utils.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}