// Copyright (c) 2024 The illium developers
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package indexers

import (
	"context"
	"encoding/binary"
	"errors"
	datastore "github.com/ipfs/go-datastore"
	"github.com/project-illium/ilxd/blockchain/pb"
	"github.com/project-illium/ilxd/repo"
	"github.com/project-illium/ilxd/types"
	"github.com/project-illium/ilxd/types/blocks"
	"github.com/project-illium/ilxd/types/transactions"
	"google.golang.org/protobuf/proto"
)

var _ Indexer = (*TxIndex)(nil)

const TxIndexName = "transaction index"

// TxIndex is and implementation of the Indexer which indexes
// transactions by their ID and maps them to the location of
// the transaction in the database. This is useful functionality
// for anyone interested in inspecting a given transaction, for
// example, block explorers.
type TxIndex struct{}

// NewTxIndex returns a new TxIndex.
func NewTxIndex() *TxIndex {
	return &TxIndex{}
}

// Key returns the key of the index as a string.
func (idx *TxIndex) Key() string {
	return repo.TxIndexKey
}

// Name returns the human-readable name of the index.
func (idx *TxIndex) Name() string {
	return TxIndexName
}

// ConnectBlock is called when a block is connected to the chain.
// The indexer can use this opportunity to parse it and store it in
// the database. The database transaction must be respected.
func (idx *TxIndex) ConnectBlock(dbtx datastore.Txn, blk *blocks.Block) error {
	for i, tx := range blk.Transactions {
		valueBytes := make([]byte, 36)
		binary.BigEndian.PutUint32(valueBytes[:4], uint32(i))
		copy(valueBytes[4:], blk.ID().Bytes())

		if err := dsPutIndexValue(dbtx, idx, tx.ID().String(), valueBytes); err != nil {
			return err
		}
	}
	if err := dsPutIndexerHeight(dbtx, idx, blk.Header.Height); err != nil {
		return err
	}
	return nil
}

// GetTransaction looks up the block id and position in the transaction index then fetches the
// transaction from the db and returns it.
func (idx *TxIndex) GetTransaction(ds repo.Datastore, txid types.ID) (*transactions.Transaction, error) {
	valueBytes, err := dsFetchIndexValue(ds, idx, txid.String())
	if err != nil {
		return nil, err
	}
	pos := binary.BigEndian.Uint32(valueBytes[:4])
	blockID := types.NewID(valueBytes[4:])

	ser, err := ds.Get(context.Background(), datastore.NewKey(repo.BlockTxsKeyPrefix+blockID.String()))
	if err != nil {
		return nil, err
	}

	var dsTxs pb.DBTxs
	if err := proto.Unmarshal(ser, &dsTxs); err != nil {
		return nil, err
	}

	if int(pos) > len(dsTxs.Transactions)-1 {
		return nil, errors.New("tx index position out of range")
	}

	return dsTxs.Transactions[pos], nil
}

// GetContainingBlockID returns the ID of the block containing the transaction.
func (idx *TxIndex) GetContainingBlockID(ds repo.Datastore, txid types.ID) (types.ID, error) {
	valueBytes, err := dsFetchIndexValue(ds, idx, txid.String())
	if err != nil {
		return types.ID{}, err
	}
	return types.NewID(valueBytes[4:]), nil
}

func (idx *TxIndex) Close(ds repo.Datastore) error {
	return nil
}

func DropTxIndex(ds repo.Datastore) error {
	return dsDropIndex(ds, &TxIndex{})
}
