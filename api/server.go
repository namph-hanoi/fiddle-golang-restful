package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/namph-hanoi/fiddle-golang-restful/db/sqlc"
	"github.com/namph-hanoi/fiddle-golang-restful/token"
	"github.com/namph-hanoi/fiddle-golang-restful/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		value.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)
	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)

	for _, route := range router.Routes() {
		fmt.Printf("%-6s %s\n", route.Method, route.Path)
	}
	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Router returns the gin engine for testing purposes
func (server *Server) Router() *gin.Engine {
	return server.router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
