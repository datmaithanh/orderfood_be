package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() *gin.Engine {
	router := gin.Default()
	
	// Public routes

	// User routes
	router.POST("/users/login", server.loginUser)

	// Customer routes
	router.POST("/customers", server.createCustomer)
	router.GET("/customers/:id", server.getCustomer)


	// Protected routes
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Auth user routes 
	authRouter.POST("/users", server.createUser)
	authRouter.GET("/users/:id", server.getUser)
	authRouter.GET("/users", server.listUsers)
	authRouter.DELETE("/users/:id", server.deleteUser)
	authRouter.PUT("users/:id", server.updateUser)
	authRouter.PATCH("users/password/:id", server.updateUserWithPassword)

	// Auth customer routes
	authRouter.GET("/customers", server.listCustomer)
	authRouter.DELETE("/customers/:id", server.deleteCustomer)
	
	// Category routes
	authRouter.POST("/categories", server.createCategory)
	authRouter.GET("/categories/:id", server.getCategory)
	authRouter.GET("/categories", server.listCategory)
	authRouter.DELETE("/categories/:id", server.deleteCategory)
	authRouter.PATCH("/categories/:id", server.updateCategory)
	server.router = router
	return router
}
