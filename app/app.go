package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quiknode-labs/paris/configs"
	"github.com/quiknode-labs/paris/controllers"
	"github.com/quiknode-labs/paris/services/zmq"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.New()
}

func Run() {
	/*
		====== Setup configs ============
	*/

	config := configs.Get()

	/*
		====== Setup datadog ============
	*/
	// stats, err := statsd.New(config.DatadogUrl)
	// if err != nil {
	// 	logrus.WithField("error", err).Error("Failed to init statsd client")
	// }

	/*
		====== Setup services ===========
	*/
	zmqServer := zmq.NewZMQService(nil)
	go zmqServer.Start(context.Background())

	/*
		====== Setup controllers ========
	*/
	rpcCtl := controllers.NewRPCController(zmqServer, config)

	/*
		====== Setup middlewares ========
	*/
	router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	/*
		====== Setup routes =============
	*/
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	router.POST("/", rpcCtl.GetRPC)

	// Run
	router.Run(":3000")
}
