package zmq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	zmq "github.com/go-zeromq/zmq4"
	gjson "github.com/goccy/go-json"
	"github.com/jellydator/ttlcache/v3"
	"github.com/quiknode-labs/paris/types"
	"github.com/sirupsen/logrus"
)

// ZMQService interface
type ZMQService interface {
	Start(context.Context) error
	GetFromLongCache(key string) []byte
	GetFromShortCache(key string) string
	processMessages(context.Context)
}

type zmqService struct {
	stats      *statsd.Client
	socket     zmq.Socket
	msgChan    chan []byte
	longCache  *ttlcache.Cache[string, []byte]
	shortCache *ttlcache.Cache[string, string]
}

// NewZMQService will instantiate ZMQ Service
func NewZMQService(stats *statsd.Client) ZMQService {
	return &zmqService{
		stats:      stats,
		socket:     zmq.NewSub(context.Background()),
		msgChan:    make(chan []byte, 1000),
		longCache:  ttlcache.New[string, []byte](ttlcache.WithTTL[string, []byte](30 * time.Minute)),
		shortCache: ttlcache.New[string, string](ttlcache.WithTTL[string, string](12 * time.Second)),
	}
}

func (n *zmqService) Start(ctx context.Context) error {
	// Connect to the endpoint
	if err := n.socket.Dial("tcp://129.213.58.232:31338"); err != nil {
		return err
	}

	// Subscribe to all topics
	if err := n.socket.SetOption(zmq.OptionSubscribe, ""); err != nil {
		return err
	}

	log.Println("subscriber connected to port 31338")

	go n.processMessages(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			//  Receive frames from socket
			msg, err := n.socket.Recv()
			if err != nil {
				logrus.WithField("error", err).Error("zmq socket recv")
				continue
			}

			// Add the message to the channel
			select {
			case n.msgChan <- msg.Frames[0]:
				// Message sent successfully
			default:
				logrus.Error("Channel is full, message not sent")
			}
		}
	}
}

func (n *zmqService) GetFromLongCache(key string) []byte {
	value := n.longCache.Get(key)
	return value.Value()
}

func (n *zmqService) GetFromShortCache(key string) string {
	value := n.shortCache.Get(key)
	return value.Value()
}

func (n *zmqService) processMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-n.msgChan:
			if len(msg) < 81 {
				log.Println("received message is too short")
				continue
			}

			topic := msg[0:3]
			blockNumber := msg[4:13]
			hash := msg[14:79]
			restOfTheData := msg[81:]
			//log.Printf("Topic: %s BlockNumber: %s Hash: %s\n", topic, blockNumber, hash)

			if bytes.Equal(topic, []byte("009")) {
				n.shortCache.Set("latest", string(restOfTheData), time.Second*12)
				continue
			}

			if bytes.Equal(topic, []byte("008")) {
				n.shortCache.Set("final", string(restOfTheData), time.Second*12)
				continue
			}

			if bytes.Equal(topic, []byte("007")) {
				n.shortCache.Set("safe", string(restOfTheData), time.Second*12)
				continue
			}

			// block ref
			blockKey := fmt.Sprintf("%s_%s", topic, blockNumber)
			n.longCache.Set(blockKey, restOfTheData, time.Hour)

			// hash ref
			hashKey := fmt.Sprintf("%s_%s", topic, hash)
			n.longCache.Set(hashKey, restOfTheData, time.Hour)

			go n.processExtraData(topic, blockNumber, hash, restOfTheData)
		}
	}
}

func unmarshalData(input []byte) ([]types.TransactionData, error) {
	var raw []map[string]interface{}
	if err := json.Unmarshal(input, &raw); err != nil {
		return nil, err
	}

	data := make([]types.TransactionData, len(raw))
	for i, item := range raw {
		transactionHash, ok := item["transactionHash"].(string)
		if !ok {
			return nil, fmt.Errorf("couldn't convert transactionHash to string")
		}

		// Remove transactionHash from the map
		delete(item, "transactionHash")

		data[i] = types.TransactionData{
			TransactionHash: transactionHash,
			Data:            item,
		}
	}

	return data, nil
}

func (n *zmqService) processExtraData(topic, blockNumber, hash, restOfTheData []byte) {
	if bytes.Equal(topic, []byte("005")) || bytes.Equal(topic, []byte("006")) {
		// debug_traceTransaction with or without log
		go func(data []byte) {
			var traces []types.TraceData
			if err := gjson.Unmarshal(data, &traces); err != nil {
				log.Printf("failed to unmarshal data: %v", err)
				return
			}

			for _, trace := range traces {
				jsonResult, err := gjson.Marshal(trace.Result)
				if err != nil {
					log.Printf("failed to marshal result: %v", err)
					continue
				}
				n.longCache.Set(trace.TxHash, jsonResult, time.Hour)
			}
		}(restOfTheData)
	}

	if bytes.Equal(topic, []byte("002")) {
		// transaction receipt
		go func(data []byte) {
			transactions, err := unmarshalData(data)
			if err != nil {
				log.Printf("failed to unmarshal data: %v", err)
				return
			}

			for _, transaction := range transactions {
				jsonData, err := gjson.Marshal(transaction.Data)
				if err != nil {
					log.Printf("failed to marshal data: %v", err)
					continue
				}
				n.longCache.Set(transaction.TransactionHash, jsonData, time.Hour)
			}
		}(restOfTheData)
	}
}
