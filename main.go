package main

import (
	"fmt"
	"github.com/lyyubava/solidgate-software-engineering-school.git/models"
	"github.com/lyyubava/solidgate-software-engineering-school.git/routers"
	"os"
)

func main() {
	models.Database(os.Getenv("DATABASE_CONNECTION_STRING"))
	routerInit := routers.InitRouter()
	routerInit.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("API_PORT")))
}
