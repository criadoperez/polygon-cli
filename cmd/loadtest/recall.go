package loadtest

import (
	"context"
	"encoding/json"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"math/big"
)

// TODO allow this to be pre-specified with an input file
func getRecentBlocks(ctx context.Context, ec *ethclient.Client, c *ethrpc.Client) ([]*json.RawMessage, error) {
	bn, err := ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	rawBlocks, err := util.GetBlockRange(ctx, bn-*inputLoadTestParams.RecallLength, bn, c)
	return rawBlocks, err
}

func getRecallTransactions(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client) ([]rpctypes.PolyTransaction, error) {
	rb, err := getRecentBlocks(ctx, c, rpc)
	if err != nil {
		return nil, err
	}
	txs := make([]rpctypes.PolyTransaction, 0)
	for _, v := range rb {
		pb := new(rpctypes.RawBlockResponse)
		err := json.Unmarshal(*v, pb)
		if err != nil {
			return nil, err
		}
		for _, t := range pb.Transactions {
			pt := rpctypes.NewPolyTransaction(&t)
			txs = append(txs, pt)
		}
	}
	return txs, nil
}

func rawTransactionToNewTx(pt rpctypes.PolyTransaction, nonce uint64, price, tipCap *big.Int) *ethtypes.Transaction {
	if pt.MaxFeePerGas() != 0 || pt.ChainID() != 0 {
		return rawTransactionToDynamicFeeTx(pt, nonce, price, tipCap)
	}
	return rawTransactionToLegacyTx(pt, nonce, price)
}
func rawTransactionToDynamicFeeTx(pt rpctypes.PolyTransaction, nonce uint64, price, tipCap *big.Int) *ethtypes.Transaction {
	toAddr := pt.To()
	chainId := new(big.Int).SetUint64(pt.ChainID())
	dynamicFeeTx := &ethtypes.DynamicFeeTx{
		ChainID:   chainId,
		To:        &toAddr,
		Data:      pt.Data(),
		Value:     pt.Value(),
		Gas:       pt.Gas(),
		GasFeeCap: price,
		GasTipCap: tipCap,
		Nonce:     nonce,
	}
	tx := ethtypes.NewTx(dynamicFeeTx)
	return tx
}
func rawTransactionToLegacyTx(pt rpctypes.PolyTransaction, nonce uint64, price *big.Int) *ethtypes.Transaction {
	toAddr := pt.To()
	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		To:       &toAddr,
		Value:    pt.Value(),
		Data:     pt.Data(),
		Gas:      pt.Gas(),
		Nonce:    nonce,
		GasPrice: price,
	})
	return tx
}
