package controllers

import (
	"encoding/json"
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

type JSONRPCResponse2 struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result"`
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

		tracer, ok := req.Params[1].(map[string]interface{})["tracer"]
		if !ok {
		} else {
			if tracer != "callTracer" {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "support only for callTracer",
				})
				return
			}
		}

		withLogs := false
		if len(req.Params) >= 2 {
			paramsMap, ok := req.Params[1].(map[string]interface{})
			if ok {
				tracerConfig, ok := paramsMap["tracerConfig"].(map[string]interface{})
				if ok {
					log, ok := tracerConfig["withLog"].(bool)
					if ok && log {
						withLogs = true
					}

					topcalls, ok := tracerConfig["onlyTopCall"].(bool)
					if ok && topcalls {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": "onlyTopCall supported as false only",
						})
						return
					}
				}
			}
		}

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
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "debug_traceTransaction":

		tracer, ok := req.Params[1].(map[string]interface{})["tracer"]
		if !ok {
		} else {
			if tracer != "callTracer" {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "support only for callTracer",
				})
				return
			}
		}

		withLogs := false
		if len(req.Params) >= 2 {
			paramsMap, ok := req.Params[1].(map[string]interface{})
			if ok {
				tracerConfig, ok := paramsMap["tracerConfig"].(map[string]interface{})
				if ok {
					log, ok := tracerConfig["withLog"].(bool)
					if ok && log {
						withLogs = true
					}

					topcalls, ok := tracerConfig["onlyTopCall"].(bool)
					if ok && topcalls {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": "onlyTopCall supported as false only",
						})
						return
					}
				}
			}
		}

		var data []byte
		if withLogs {
			data = ctl.zmqServer.GetFromLongCache("005_" + req.Params[0].(string))
		} else {
			data = ctl.zmqServer.GetFromLongCache("006_" + req.Params[0].(string))
		}

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "eth_getTransactionReceipt":
		data := ctl.zmqServer.GetFromLongCache("002_" + req.Params[0].(string))

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "eth_getBlockReceipts":
		data := ctl.zmqServer.GetFromLongCache("002_" + ctl.translateToBlockNumber(req.Params[0].(string)))

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "eth_getBlockByNumber":
		if len(req.Params) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "wrong number of parameters",
			})
			return
		}

		fulltx := req.Params[1].(bool)
		var data []byte
		if fulltx {
			data = ctl.zmqServer.GetFromLongCache("001_" + ctl.translateToBlockNumber(req.Params[0].(string)))
		} else {
			data = ctl.zmqServer.GetFromLongCache("000_" + ctl.translateToBlockNumber(req.Params[0].(string)))
		}

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return
	case "eth_getBlockByHash":
		if len(req.Params) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "wrong number of parameters",
			})
			return
		}

		fulltx := req.Params[1].(bool)
		var data []byte
		if fulltx {
			data = ctl.zmqServer.GetFromLongCache("001_" + ctl.translateToBlockNumber(req.Params[0].(string)))
		} else {
			data = ctl.zmqServer.GetFromLongCache("000_" + ctl.translateToBlockNumber(req.Params[0].(string)))
		}

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return
	case "eth_getBalance":
		data := ctl.zmqServer.GetFromLongCache("003_" + ctl.translateToBlockNumber(req.Params[0].(string)))

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return
	case "eth_getCode":
		data := ctl.zmqServer.GetFromLongCache("004_" + ctl.translateToBlockNumber(req.Params[0].(string)))

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fullfil",
			})
			return
		}

		response := JSONRPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "eth_gasPrice":
		data := ctl.zmqServer.GetFromShortCache("gasPrice")

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fulfill",
			})
			return
		}

		response := JSONRPCResponse2{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return

	case "eth_blockNumber":
		data := ctl.zmqServer.GetFromShortCache("latest")

		if len(data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no data in the cache to fulfill",
			})
			return
		}

		response := JSONRPCResponse2{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  data,
		}

		c.JSON(http.StatusOK, response)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown method",
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "cant fullfil",
	})
}
