# `polycli loadtest`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Run a generic load test against an Eth/EVM style JSON-RPC endpoint.

```bash
polycli loadtest url [flags]
```

## Usage

The `loadtest` tool is meant to generate various types of load against RPC end points. It leverages the [`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) library Go Ethereum to interact with the blockchain.x

```bash
$ polycli wallet inspect  --mnemonic "code code code code code code code code code code code quality" --addresses 1
```

The `--mode` flag is important for this command.

- `t` will only perform transfers to the `--to-address`. This is a fast and common operation.
- `d` will deploy the load testing contract over and over again.
- `c` will call random functions in our load test contract.
- `f` will call a specific function on the load test contract. The function is specified using the `-f` flag
- `2` will run an ERC20 transfer test. It starts out by minting a large amount of an ERC20 contract then transferring it in small amounts.
- `7` will run an ERC721 test which will mint an NFT over and over again.
- `i` will call the increment function repeatedly on the load test contract. It's a minimal example of a contract call that will require an update to a contract's storage.
- `r` will call any of th eother modes randomly.
- `s` is used for Avail / Eth to store random data in large amounts.
- `l` will call a smart contract function that runs as long as it can, based on the block limit.

The default private key is: `42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa`. We can use `wallet inspect` to get more information about this address, in particular its `ETHAddress` if you want to check balance or pre-mine value for this particular account.

Here is a simple example that runs 1000 requests at a max rate of 1 request per second against the http rpc endpoint on localhost. It's running in transaction mode so it will perform simple transactions send to the default address.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 1 --mode t http://localhost:8888
```

Another example, a bit slower, and that specifically calls the [LOG4](https://www.evm.codes/#a4) function in the load test contract in a loop for 25,078 iterations. That number was picked specifically to require almost all of the gas for a single transaction.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 50 --rate-limit 0.5  --mode f --function 164 --iterations 25078 http://private.validator-001.devnet02.pos-v3.polygon.private:8545
```

### Load Test Contract

The codebase has a contract that used for load testing. It's written in Yul and Solidity. The workflow for modifying this contract is.

1. Make changes to <file:contracts/LoadTester.sol>
2. Compile the contracts:
   - `$ solc LoadTester.sol --bin --abi -o . --overwrite`
3. Run `abigen`
   - `$ abigen --abi LoadTester.abi --pkg contracts --type LoadTester --bin LoadTester.bin --out loadtester.go`
4. Run the loadtester to enure it deploys and runs successfully
   - `$ polycli loadtest --verbosity 700 http://127.0.0.1:8541`

### Avail / Substrate

The loadtest tool works with Avail, but not with the same level of functionality. There's no EVM so the functional calls will not work. This is a basic example which would transfer value in a loop 10 times.

```bash
$ polycli loadtest --app-id 0 --to-random=true  --data-avail --verbosity 700 --chain-id 42 --concurrency 1 --requests 10 --rate-limit 1 --mode t 'http://devnet01.dataavailability.link:8545'
```

This is a similar test but storing random nonsense hexwords.

```bash
$ polycli loadtest --app-id 0 --data-avail --verbosity 700 --chain-id 42 --concurrency 1 --requests 10 --rate-limit 1 --mode s --byte-count 16384 'http://devnet01.dataavailability.link:8545'
```

## Flags

```bash
      --adaptive-backoff-factor float              When using adaptive rate limiting, this flag controls our multiplicative decrease value. (default 2)
      --adaptive-cycle-duration-seconds uint       When using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates (default 10)
      --adaptive-rate-limit                        Enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint         When using adaptive rate limiting, this flag controls the size of the additive increases. (default 50)
      --batch-size uint                            Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time. (default 999)
  -b, --byte-count uint                            If we're in store mode, this controls how many bytes we'll try to store in our contract (default 1024)
      --call-only                                  When using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features.
      --chain-id uint                              The chain id for the transactions.
  -c, --concurrency int                            Number of requests to perform concurrently. Default is one request at a time. (default 1)
      --contract-call-block-interval uint          During deployment, this flag controls if we should check every block, every other block, or every nth block to determine that the contract has been deployed (default 1)
      --contract-call-nb-blocks-to-wait-for uint   The number of blocks to wait for before giving up on a contract deployment (default 30)
      --data-avail                                 [DEPRECATED] Enables Avail load testing
      --erc20-address string                       The address of a pre-deployed erc 20 contract
      --erc721-address string                      The address of a pre-deployed erc 721 contract
      --force-contract-deploy                      Some load test modes don't require a contract deployment. Set this flag to true to force contract deployments. This will still respect the --lt-address flags.
  -f, --function --mode f                          A specific function to be called if running with --mode f or a specific precompiled contract when running with `--mode a` (default 1)
      --gas-limit uint                             In environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas
      --gas-price uint                             In environments where the gas price can't be determined automatically, we can specify it manually
  -h, --help                                       help for loadtest
  -i, --iterations uint                            If we're making contract calls, this controls how many times the contract will execute the instruction in a loop. If we are making ERC721 Mints, this indicates the minting batch size (default 1)
      --legacy                                     Send a legacy transaction instead of an EIP1559 transaction.
      --lt-address string                          The address of a pre-deployed load test contract
  -m, --mode string                                The testing mode to use. It can be multiple like: "tcdf"
                                                   t - sending transactions
                                                   d - deploy contract
                                                   c - call random contract functions
                                                   f - call specific contract function
                                                   p - call random precompiled contracts
                                                   a - call a specific precompiled contract address
                                                   s - store mode
                                                   r - random modes
                                                   2 - ERC20 Transfers
                                                   7 - ERC721 Mints
                                                   R - total recall (default "t")
      --output-mode string                         Format mode for summary output (json | text) (default "text")
      --priority-gas-price uint                    Specify Gas Tip Price in the case of EIP-1559
      --private-key string                         The hex encoded private key that we'll use to send transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --rate-limit float                           An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together (default 4)
      --recall-blocks uint                         The number of blocks that we'll attempt to fetch for recall (default 50)
  -n, --requests int                               Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results. (default 1)
      --seed int                                   A seed for generating random values and addresses (default 123456)
      --send-amount string                         The amount of wei that we'll send every transaction (default "0x38D7EA4C68000")
      --steady-state-tx-pool-size uint             When using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off. (default 1000)
      --summarize                                  Should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time
  -t, --time-limit int                             Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit. (default -1)
      --to-address string                          The address that we're going to send to (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
      --to-random                                  When doing a transfer test, should we send to random addresses rather than DEADBEEFx5
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
