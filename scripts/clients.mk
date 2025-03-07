##@ Clients
PORT?=8545

.PHONY: geth
geth: ## Start a local geth node.
	geth --dev --dev.period 2 --http --http.addr localhost --http.port $(PORT) --http.api admin,debug,web3,eth,txpool,personal,miner,net --verbosity 5 --rpc.gascap 50000000  --rpc.txfeecap 0 --miner.gaslimit  10 --miner.gasprice 1 --gpo.blocks 1 --gpo.percentile 1 --gpo.maxprice 10 --gpo.ignoreprice 2 --dev.gaslimit 50000000

.PHONY: avail
avail: ## Start a local avail node.
	avail --dev --rpc-port $(PORT)

LOADTEST_ACCOUNT=0x85da99c8a7c2c95964c8efd687e95e632fc533d6
LOADTEST_FUNDING_AMOUNT_ETH=5000
eth_coinbase := $(shell curl -s -H 'Content-Type: application/json' -d '{"jsonrpc": "2.0", "id": 2, "method": "eth_coinbase", "params": []}' http://127.0.0.1:${PORT} | jq -r ".result")
hex_funding_amount := $(shell echo "obase=16; ${LOADTEST_FUNDING_AMOUNT_ETH}*10^18" | bc)
.PHONY: geth-loadtest
geth-loadtest: build ## Fund test account with 5k ETH and run loadtest against an EVM/Geth chain.
	curl -H "Content-Type: application/json" -d '{"jsonrpc":"2.0", "method":"eth_sendTransaction", "params":[{"from": "${eth_coinbase}","to": "${LOADTEST_ACCOUNT}","value": "0x${hex_funding_amount}"}], "id":1}' http://127.0.0.1:${PORT}
	sleep 5
	$(BUILD_DIR)/$(BIN_NAME) loadtest --verbosity 700 --chain-id 1337 --concurrency 1 --requests 1000 --rate-limit 100 --mode c --legacy http://127.0.0.1:$(PORT)

.PHONY: avail-loadtest
avail-loadtest: build ## Run loadtest against an Avail chain.
	$(BUILD_DIR)/$(BIN_NAME) loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 5 --mode t --data-avail http://127.0.0.1:$(PORT)
