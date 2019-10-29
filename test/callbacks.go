package test

import (
	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	"github.com/filecoin-project/go-http-api/types"
	"github.com/ipfs/go-cid"
	"math/big"
)

// HappyPathCallbacks provides a  dummy set of callbacks (incomplete) for testing
var HappyPathCallbacks = v1.Callbacks{
	GetActorByID:  gabid,
	GetActorNonce: gan,
	GetActors:     ga,
	GetBlockByID:  gbid,
	CreateMessage: cm,
}

func gabid(addr string) (*types.Actor, error) {
	return &types.Actor{
		ActorType: "account",
		Address:   addr,
		Code:      cid.Cid{},
		Nonce:     big.NewInt(100),
		Balance:   big.NewInt(1000),
		Exports:   map[string]types.ReadableFunctionSignature{},
		Head:      cid.Cid{},
	}, nil
}
func gan(_ string) (*big.Int, error) {
	return big.NewInt(1234), nil
}
func ga() ([]*types.Actor, error) {
	a, err := gabid("foo")
	return []*types.Actor{a}, err
}

func gbid(bid string) (*types.Block, error) {
	bh := types.BlockHeader{
		Miner:                 fixtures.TestAddress0,
		Height:                234,
		ParentStateRoot:       cid.Cid{},
		ParentMessageReceipts: cid.Cid{},
		Messages:              cid.Cid{},
		Timestamp:             3984954,
		BlockSig:              []byte("sdlkfsdlfkjsdflkjfs"),
	}
	bcid, err := cid.Parse(bid)
	if err != nil {
		return &types.Block{}, err
	}
	return &types.Block{
		ID:     bcid,
		Header: bh,
	}, nil
}

func cm(m *types.Message) (*types.SignedMessage, error) {
	return &types.SignedMessage{
		Value:      m.Value,
		GasPrice:   m.GasPrice,
		GasLimit:   m.GasLimit,
		Method:     m.Method,
		Parameters: m.Parameters,
	}, nil
}
