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
	
	//Auth Category routes
	authRouter.POST("/categories", server.createCategory)
	authRouter.GET("/categories/:id", server.getCategory)
	authRouter.GET("/categories", server.listCategory)
	authRouter.DELETE("/categories/:id", server.deleteCategory)
	authRouter.PATCH("/categories/:id", server.updateCategory)

	//Auth Menu routes
	authRouter.POST("/menus", server.createMenu)
	authRouter.GET("/menus/:id", server.getMenu)
	authRouter.GET("/menus", server.listMenu)
	authRouter.DELETE("/menus/:id", server.deleteMenu)
	authRouter.PATCH("/menus/:id", server.updateMenu)

	// Auth Table routes
	authRouter.POST("/tables", server.createTable)
	authRouter.GET("/tables/:id", server.getTable)
	authRouter.GET("/tables", server.listTables)
	authRouter.PATCH("/tables/:id", server.updateTableStatus)
	authRouter.DELETE("/tables/:id", server.deleteTable)

	// Auth Order routes
	authRouter.POST("/orders", server.createOrder)
	authRouter.GET("/orders/:id", server.getOrder)
	authRouter.GET("/orders", server.listOrders)
	authRouter.DELETE("/orders/:id", server.deleteOrder)
	authRouter.PUT("/orders/:id", server.updateOrder)
	authRouter.PATCH("/orders/status/:id", server.updateOrderStatus)

	// Auth Order Item routes
	authRouter.POST("/orderitems", server.createOrderItem)
	authRouter.GET("/orderitems/:id", server.getOrderItem)
	authRouter.GET("/orderitems", server.listOrderItems)
	authRouter.DELETE("/orderitems/:id", server.deleteOrderItem)
	authRouter.PUT("/order_items/:id", server.updateOrderItem)

	// Auth Payment routes
	authRouter.POST("/payments", server.createPayment)
	authRouter.GET("/payments/:id", server.getPayment)
	authRouter.GET("/payments", server.listPayments)
	authRouter.DELETE("/payments/:id", server.deletePayment)
	authRouter.PATCH("/payments/status/:id", server.updatePaymentStatus)


	server.router = router
	return router
}
