package main

import (
	"data-management/src/config"
	"data-management/src/controller"
	"data-management/src/middleware"
	"fmt"
	"log"
	"os"
	"time"

	"data-management/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	"github.com/subosito/gotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	if err := gotenv.Load(); err != nil {
		log.Println(aurora.Red(err))
	}

	fmt.Println("Using timezone:", aurora.Green(time.Now().Location().String()))
}

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description E.g. Bearer Your.Token
func main() {
	appPort := ":" + os.Getenv(config.ENV_KEY_PORT)

	if os.Getenv("LOCAL") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(cors.Default())

	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiV1 := router.Group(basePath)
	apiV1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Nuwhofev"}) })

	apiV1Analyze := router.Group(basePath + "/analytics")
	apiV1Analyze.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Nuwhofev"}) })

	jwtMiddleware, err := middleware.InitJwt()
	if err != nil {
		log.Println(err)
		return
	}
	controller.NewPublicController(apiV1, jwtMiddleware)
	// apiV1.Use(jwtMiddleware.MiddlewareFunc())
	// apiV1Analyze.Use(jwtMiddleware.MiddlewareFunc())

	//* ------------------------ REGISTER CRUD CONTROLLER ------------------------ */
	//* CRUD API */

	controller.NewUserController(apiV1)
	controller.NewRoleController(apiV1)
	controller.NewDataSetController(apiV1)
	controller.NewCatalogController(apiV1)
	controller.NewBlogController(apiV1)

	log.Println(aurora.Green(
		fmt.Sprintf("http://localhost%s/swagger/index.html", appPort),
	))
	log.Fatalln(router.Run(appPort))
}
