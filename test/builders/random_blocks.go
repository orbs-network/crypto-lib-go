// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package builders

import (
	"github.com/orbs-network/crypto-lib-go/test/rand"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

func RandomizedBlockChain(numBlocks int32, ctrlRand *rand.ControlledRand) []*protocol.BlockPairContainer {
	return RandomizedBlockChainWithLimit(numBlocks, ctrlRand, 200, 200)
}

func RandomizedBlockChainWithLimit(numBlocks int32, ctrlRand *rand.ControlledRand, maxTransactions uint32, maxStateDiffs uint32) []*protocol.BlockPairContainer {
	blocks := make([]*protocol.BlockPairContainer, 0, numBlocks)

	var prev *protocol.BlockPairContainer
	for bi := 1; bi <= cap(blocks); bi++ {
		newBlock := RandomizedBlockWithLimit(primitives.BlockHeight(bi), ctrlRand, prev, maxTransactions, maxStateDiffs)
		blocks = append(blocks, newBlock)
		prev = newBlock
	}
	return blocks
}

func RandomizedBlock(h primitives.BlockHeight, ctrlRand *rand.ControlledRand, prev *protocol.BlockPairContainer) *protocol.BlockPairContainer {
	return RandomizedBlockWithLimit(h, ctrlRand, prev, 200, 200)
}

func RandomizedBlockWithLimit(h primitives.BlockHeight, ctrlRand *rand.ControlledRand, prev *protocol.BlockPairContainer, maxTransactions uint32, maxStateDiffs uint32) *protocol.BlockPairContainer {
	builder := BlockPair().
		WithHeight(h).
		WithTransactions(ctrlRand.Uint32() % maxTransactions).
		WithStateDiffs(ctrlRand.Uint32() % maxStateDiffs).
		WithReceiptsForTransactions().
		WithEmptyLeanHelixBlockProof()
	if prev != nil {
		builder.WithPrevBlock(prev)
	}
	return builder.Build()
}
