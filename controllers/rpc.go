package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	gjson "github.com/goccy/go-json"
	"github.com/quiknode-labs/paris/configs"
	"github.com/quiknode-labs/paris/services/zmq"
)

// AppController interface
type RPCController interface {
	GetRPC(*gin.Context)
}

type rpcController struct {
	zmqServer zmq.ZMQService
	config    *configs.AppConfig
}

func NewRPCController(zmqServer zmq.ZMQService, c *configs.AppConfig) RPCController {
	return &rpcController{
		zmqServer: zmqServer,
		config:    c,
	}
}

type RPCRequest struct {
	Method string        `json:"method"`
	ID     interface{}   `json:"id"`
	Params []interface{} `json:"params"`
}

type JSONRPCResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result"`
}

func (ctl *rpcController) translateToBlockNumber(param1 string) string {
	if param1 == "latest" {
		return string(ctl.zmqServer.GetFromShortCache("latest"))
	} else if param1 == "safe" {
		return string(ctl.zmqServer.GetFromShortCache("safe"))
	} else if param1 == "final" {
		return string(ctl.zmqServer.GetFromShortCache("final"))
	} else {
		return param1
	}
}

func (ctl *rpcController) GetRPC(c *gin.Context) {
	var req RPCRequest
	err := gjson.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	switch req.Method {
	case "debug_traceBlockByNumber":

		withLogs := false
		if len(req.Params) >= 2 {
			log := req.Params[1].(map[string]interface{})["tracerConfig"].(map[string]interface{})["withLog"]
			if log != nil {
				if log == "true" {
					withLogs = true
				}
			}
		}

		fmt.Println("EAAAAAAAAAAAAAAAA", req.Params[0].(string))

		var data []byte
		if withLogs {
			data = ctl.zmqServer.GetFromLongCache("005_" + ctl.translateToBlockNumber(req.Params[0].(string)))
		} else {
			data = ctl.zmqServer.GetFromLongCache("006_" + ctl.translateToBlockNumber(req.Params[0].(string)))
		}

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "debug_traceTransaction":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	case "eth_getTransactionReceipt":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	case "eth_getBlockReceipts":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	case "eth_getBlockByNumber":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	case "eth_getBlockByHash":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	case "eth_getBalance":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	case "eth_getCode":
		// handle method 2
		// ...
		c.JSON(http.StatusOK, gin.H{
			"message": "Method 2 executed",
			// ... return any other necessary info
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown method",
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "cant fullfil",
	})
}
