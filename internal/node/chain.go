package node

import (
	"encoding/hex"
	"fmt"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/types"
)

type Chain struct {
	blockStore   BlockStorer
	blockHeaders *BlockHeaderList
}

func NewChain(bs BlockStorer) *Chain {
	return &Chain{
		blockStore:   bs,
		blockHeaders: NewBlockHeaderList(),
	}
}

func (c *Chain) AddBlock(b *genproto.Block) error {
	c.blockHeaders.Add(b.Header)
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHash(hash []byte) (*genproto.Block, error) {
	hashString := hex.EncodeToString(hash)
	return c.blockStore.Get(hashString)
}

func (c *Chain) GetBlockByHeight(height int) (*genproto.Block, error) {
	if height > c.Height() {
		return nil, fmt.Errorf("block with height %d doesn't exist", height)
	}

	blockHeader := c.blockHeaders.Get(height)
	hash := types.CalcBlockHeaderHash(blockHeader)
	return c.GetBlockByHash(hash)
}

func (c *Chain) Height() int {
	return c.blockHeaders.Height()
}

//-----------------------------------------------------------------------------
//  BlockHeaderList
//-----------------------------------------------------------------------------

type BlockHeaderList struct {
	headerList []*genproto.BlockHeader
}

func NewBlockHeaderList() *BlockHeaderList {
	return &BlockHeaderList{
		headerList: make([]*genproto.BlockHeader, 0),
	}
}

func (hs *BlockHeaderList) Add(h *genproto.BlockHeader) {
	hs.headerList = append(hs.headerList, h)
}

func (hs *BlockHeaderList) Get(height int) *genproto.BlockHeader {
	return hs.headerList[height]
}

func (hs *BlockHeaderList) Height() int {
	// blockchain always has a genesis block
	return len(hs.headerList) - 1
}

func (hs *BlockHeaderList) Length() int {
	return len(hs.headerList) - 1
}
