package api

import "github.com/gin-gonic/gin"

func(server *Server) setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/users", server.createUser)

	server.router = router
	return router
}
