# paris

```
ab -p postdata.json -T 'application/json' -c 1000 -n 100000 http://localhost:3000/

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","id":1,"method":"debug_traceBlockByNumber","params":["0x10f4989",{"tracer":"callTracer","tracerConfig":{"withLog":true,"onlyTopCall":false}}]}' | jq | more

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","id":1,"method":"debug_traceTransaction","params":["0x4d492dc0c6aee0a6bade404373b23a83bc20e3763cb131bbbab57cd95e799437",{"tracer":"callTracer","tracerConfig":{"withLog":true,"onlyTopCall":false}}]}' | jq

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"method":"eth_getTransactionReceipt","params":["0x7ea5554f05305a4eeaeec394504ffa76cc87ff93f95d39c150b0d9605fe71af6"],"id":1,"jsonrpc":"2.0"}' | jq

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"method":"eth_getBlockReceipts","params":["latest"],"id":1,"jsonrpc":"2.0"}' | jq

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"method":"eth_getBlockByNumber","params":["latest",true],"id":1,"jsonrpc":"2.0"}' | jq

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"method":"eth_getBlockByHash","params":["latest",false],"id":1,"jsonrpc":"2.0"}' | jq

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"method":"eth_gasPrice","params":[],"id":1,"jsonrpc":"2.0"}' | jq

curl -s localhost:3000 -X POST -H "Content-Type: application/json" --data '{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}' | jq

```


  paris:
    url: http://127.0.0.1:3000
methods:
  debug_traceBlockByNumber: paris
  debug_traceTransaction: paris
  eth_getTransactionReceipt: paris
  eth_getBlockReceipts: paris
  eth_getBlockByNumber: paris
  eth_getBlockByHash: paris
  eth_getBalance: paris
  eth_getCode: paris
  eth_gasPrice: paris


