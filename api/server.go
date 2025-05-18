package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/namph-hanoi/fiddle-golang-restful/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	for _, route := range router.Routes() {
		fmt.Printf("%-6s %s\n", route.Method, route.Path)
	}
	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
