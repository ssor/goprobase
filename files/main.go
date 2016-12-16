package main

import (
	"github.com/gin-gonic/gin"

	"flag"
)

var (
	configFile    = flag.String("config", "conf/config.json", "config file for system")
	listeningPort = flag.String("port", "8080", "listeningPort")

)

func main() {

	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	router := gin.Default()
	router.Static("/javascripts", "static/js")
	// router.Static("/images", "static/img")
	// router.Static("/stylesheets", "static/css")

	router.LoadHTMLGlob("views/*.tpl")

	router.Run(":" + *listeningPort)
}
