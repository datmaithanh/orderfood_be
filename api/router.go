package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() *gin.Engine {
	router := gin.Default()
	
	router.POST("/users/login", server.loginUser)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.POST("/users", server.createUser)
	authRouter.GET("/users/:id", server.getUser)
	authRouter.GET("/users", server.listUsers)
	authRouter.DELETE("/users/:id", server.deleteUser)
	authRouter.PUT("users/:id", server.updateUser)
	authRouter.PATCH("users/password/:id", server.updateUserWithPassword)

	server.router = router
	return router
}
