package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() *gin.Engine {
	router := gin.Default()
	

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUsers)
	router.DELETE("/users/:id", server.deleteUser)
	router.PUT("users/:id", server.updateUser)
	router.PATCH("users/password/:id", server.updateUserWithPassword)

	server.router = router
	return router
}
