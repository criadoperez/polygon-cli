package p2p

import (
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// MessageCount is used to help the outer goroutine to receive summary of the
// number and type of messages that were sent. This is used for distributed
// logging. It can be used to count the different types of messages received
// across all peer connections to provide a summary.
type MessageCount struct {
	BlockHeaders        int32 `json:",omitempty"`
	BlockBodies         int32 `json:",omitempty"`
	Blocks              int32 `json:",omitempty"`
	BlockHashes         int32 `json:",omitempty"`
	BlockHeaderRequests int32 `json:",omitempty"`
	BlockBodiesRequests int32 `json:",omitempty"`
	Transactions        int32 `json:",omitempty"`
	TransactionHashes   int32 `json:",omitempty"`
	TransactionRequests int32 `json:",omitempty"`
	Pings               int32 `json:",omitempty"`
	Errors              int32 `json:",omitempty"`
	Disconnects         int32 `json:",omitempty"`
}

// LogMessageCount will log the message counts and reset them based on the
// ticker. This should be called as with a goroutine or it will block
// indefinitely.
func LogMessageCount(count *MessageCount, ticker *time.Ticker) {
	for {
		if _, ok := <-ticker.C; !ok {
			return
		}

		c := MessageCount{
			BlockHeaders:        atomic.LoadInt32(&count.BlockHeaders),
			BlockBodies:         atomic.LoadInt32(&count.BlockBodies),
			Blocks:              atomic.LoadInt32(&count.Blocks),
			BlockHashes:         atomic.LoadInt32(&count.BlockHashes),
			BlockHeaderRequests: atomic.LoadInt32(&count.BlockHeaderRequests),
			BlockBodiesRequests: atomic.LoadInt32(&count.BlockBodiesRequests),
			Transactions:        atomic.LoadInt32(&count.Transactions),
			TransactionHashes:   atomic.LoadInt32(&count.TransactionHashes),
			TransactionRequests: atomic.LoadInt32(&count.TransactionRequests),
			Pings:               atomic.LoadInt32(&count.Pings),
			Errors:              atomic.LoadInt32(&count.Errors),
			Disconnects:         atomic.LoadInt32(&count.Disconnects),
		}

		if c.BlockHeaders+c.BlockBodies+c.Blocks+c.BlockHashes+
			c.BlockHeaderRequests+c.BlockBodiesRequests+c.Transactions+
			c.TransactionHashes+c.TransactionRequests+c.Pings+c.Errors+
			c.Disconnects == 0 {
			continue
		}

		log.Info().Interface("counts", c).Msg("Received messages")

		atomic.StoreInt32(&count.BlockHeaders, 0)
		atomic.StoreInt32(&count.BlockBodies, 0)
		atomic.StoreInt32(&count.Blocks, 0)
		atomic.StoreInt32(&count.BlockHashes, 0)
		atomic.StoreInt32(&count.BlockHeaderRequests, 0)
		atomic.StoreInt32(&count.BlockBodiesRequests, 0)
		atomic.StoreInt32(&count.Transactions, 0)
		atomic.StoreInt32(&count.TransactionHashes, 0)
		atomic.StoreInt32(&count.TransactionRequests, 0)
		atomic.StoreInt32(&count.Pings, 0)
		atomic.StoreInt32(&count.Errors, 0)
		atomic.StoreInt32(&count.Disconnects, 0)
	}
}
