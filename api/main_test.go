package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// comment the line to get the log into
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
