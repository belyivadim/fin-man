package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
  swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/belyivadim/fin-man/docs"
	"github.com/belyivadim/fin-man/controllers"
	"github.com/belyivadim/fin-man/database"
	"github.com/belyivadim/fin-man/env"
	"github.com/belyivadim/fin-man/models"
)

func main() {
  r := gin.Default()

  dns, err := env.GetDbDns()
  if err != nil {
    log.Fatal("Cannot get database dns from environment\n")
  }

  db, err := database.NewGormDB(dns)
  if err != nil {
    log.Fatalf("Cannot connect to the database with dns: '%s'\n", dns)
  }

  AutoMigrateDb(db)
  setupHandlers(r, db)

  r.Run(":8080")

}

func AutoMigrateDb(db *gorm.DB) {
  db.AutoMigrate(&models.User{})
  db.AutoMigrate(&models.Category{})
}

func setupHandlers(r *gin.Engine, db *gorm.DB) {
  r.GET("/healthcheck", func(c *gin.Context) {
    c.JSON(http.StatusOK, "Server is running...")
  })

  docs.SwaggerInfo.BasePath = "/api/v1"
  v1 := r.Group("/api/v1")

  auth := v1.Group("/auth", func(c *gin.Context) {
    c.Set("AuthService", database.NewAuthGormService(db))
    c.Next()
  })

  {
    auth.POST("/signup", controllers.Signup)
    auth.POST("/signin", controllers.Signin)

    // debug
    auth.GET("/users", controllers.GetAllUsers)
  }

  categories := v1.Group("/categories", func(c *gin.Context) {
    c.Set("CategoryService", database.NewCategoryGormService(db))
    c.Next()
  })
  categories.GET("", controllers.GetAllCategories)
  categories.GET(":id", controllers.GetCategoryById)
  categories.POST("", controllers.AddCategory)
  categories.PUT(":id", controllers.UpdateCategory)
  categories.DELETE(":id", controllers.RemoveCategory)

  expenses := v1.Group("/expenses")
  expenses.GET(":id", controllers.GetExpenseById)
  expenses.POST("", controllers.AddExpense)
  expenses.DELETE(":id", controllers.RemoveExpense)

  incomes := v1.Group("/incomes")
  incomes.GET(":id", controllers.GetIncomeById)
  incomes.POST("", controllers.AddIncome)
  incomes.DELETE(":id", controllers.RemoveIncome)

  savings := v1.Group("/savings")
  savings.GET("", controllers.GetSavingGoal)
  savings.POST("", controllers.AddSavingGoal)
  savings.PUT("", controllers.UpdateSavingGoal)
  savings.DELETE("", controllers.RemoveSavingGoal)

  v1.GET("/reports", controllers.GenerateReportForPeriod)

  r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
