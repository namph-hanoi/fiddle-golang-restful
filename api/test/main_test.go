package api_test

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/namph-hanoi/fiddle-golang-restful/api"
	db "github.com/namph-hanoi/fiddle-golang-restful/db/sqlc"
	"github.com/namph-hanoi/fiddle-golang-restful/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *api.Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := api.NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	// comment the line to get the log into
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
