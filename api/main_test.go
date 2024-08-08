package api

import (
	"os"
	"testing"
	"time"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Configuration{
		TokenSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}