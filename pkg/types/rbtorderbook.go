package types

import (
	"fmt"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/pkg/errors"
)

//go:generate callbackgen -type RBTOrderBook
type RBTOrderBook struct {
	Symbol string
	Bids   *RBTree
	Asks   *RBTree

	loadCallbacks   []func(book *RBTOrderBook)
	updateCallbacks []func(book *RBTOrderBook)
}

func NewRBOrderBook(symbol string) *RBTOrderBook {
	return &RBTOrderBook{
		Symbol: symbol,
		Bids:   NewRBTree(),
		Asks:   NewRBTree(),
	}
}

func (b *RBTOrderBook) BestBid() (PriceVolume, bool) {
	right := b.Bids.Rightmost(b.Bids.Root)
	if right != nil {
		return PriceVolume{Price: right.Key, Volume: right.Value}, true
	}
	return PriceVolume{}, false
}

func (b *RBTOrderBook) BestAsk() (PriceVolume, bool) {
	left := b.Asks.Leftmost(b.Bids.Root)
	if left != nil {
		return PriceVolume{Price: left.Key, Volume: left.Value}, true
	}
	return PriceVolume{}, false
}

func (b *RBTOrderBook) Spread() (fixedpoint.Value, bool) {
	bestBid, ok := b.BestBid()
	if !ok {
		return 0, false
	}

	bestAsk, ok := b.BestAsk()
	if !ok {
		return 0, false
	}

	return bestAsk.Price - bestBid.Price, true
}

func (b *RBTOrderBook) IsValid() (bool, error) {
	bid, hasBid := b.BestBid()
	ask, hasAsk := b.BestAsk()

	if !hasBid {
		return false, errors.New("empty bids")
	}

	if !hasAsk {
		return false, errors.New("empty asks")
	}

	if bid.Price > ask.Price {
		return false, fmt.Errorf("bid price %f > ask price %f", bid.Price.Float64(), ask.Price.Float64())
	}

	return true, nil
}

func (b *RBTOrderBook) Load(book SliceOrderBook) {
	b.Reset()
	b.update(book)
	b.EmitLoad(b)
}

func (b *RBTOrderBook) Update(book SliceOrderBook) {
	b.update(book)
	b.EmitUpdate(b)
}

func (b *RBTOrderBook) Reset() {
	b.Bids = NewRBTree()
	b.Asks = NewRBTree()
}

func (b *RBTOrderBook) updateAsks(pvs PriceVolumeSlice) {
	for _, pv := range pvs {
		if pv.Volume == 0 {
			b.Asks.Delete(pv.Price)
		} else {
			b.Asks.Upsert(pv.Price, pv.Volume)
		}
	}
}

func (b *RBTOrderBook) updateBids(pvs PriceVolumeSlice) {
	for _, pv := range pvs {
		if pv.Volume == 0 {
			b.Bids.Delete(pv.Price)
		} else {
			b.Bids.Upsert(pv.Price, pv.Volume)
		}
	}
}

func (b *RBTOrderBook) update(book SliceOrderBook) {
	b.updateBids(book.Bids)
	b.updateAsks(book.Asks)
}

func (b *RBTOrderBook) load(book SliceOrderBook) {
	b.Reset()
	b.updateBids(book.Bids)
	b.updateAsks(book.Asks)
}

func (b *RBTOrderBook) Print() {
	b.Bids.PostorderOf(b.Bids.Root, func(n *RBNode) {
		fmt.Printf("bid: %f x %f", n.Key.Float64(), n.Value.Float64())
	})

	b.Asks.PostorderOf(b.Asks.Root, func(n *RBNode) {
		fmt.Printf("ask: %f x %f", n.Key.Float64(), n.Value.Float64())
	})
}
