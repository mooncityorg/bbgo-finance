package xpuremaker

import (
	"context"
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/types"
)

func init() {
	bbgo.RegisterExchangeStrategy("xpuremaker", &Strategy{})
}

type Strategy struct {
	Symbol       string           `json:"symbol"`
	Side         string           `json:"side"`
	NumOrders    int              `json:"numOrders"`
	BehindVolume fixedpoint.Value `json:"behindVolume"`
	PriceTick    fixedpoint.Value `json:"priceTick"`
	BaseQuantity fixedpoint.Value `json:"baseQuantity"`
	BuySellRatio float64          `json:"buySellRatio"`

	book         *types.StreamOrderBook
	activeOrders map[string]types.Order
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {
	session.Subscribe(types.BookChannel, s.Symbol, types.SubscribeOptions{})

	s.book = types.NewStreamBook(s.Symbol)
	s.book.BindStream(session.Stream)

	s.activeOrders = make(map[string]types.Order)

	// We can move the go routine to the parent level.
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		s.update(orderExecutor, session)

		for {
			select {
			case <-ctx.Done():
				return

			case <-s.book.C:
				s.update(orderExecutor, session)

			case <-ticker.C:
				s.update(orderExecutor, session)
			}
		}
	}()

	return nil
}

func (s *Strategy) cancelOrders(session *bbgo.ExchangeSession) {
	var deletedIDs []string
	for clientOrderID, o := range s.activeOrders {
		log.Infof("canceling order: %+v", o)

		if err := session.Exchange.CancelOrders(context.Background(), o); err != nil {
			log.WithError(err).Error("cancel order error")
			continue
		}

		deletedIDs = append(deletedIDs, clientOrderID)
	}
	s.book.C.Drain(1*time.Second, 3*time.Second)

	for _, id := range deletedIDs {
		delete(s.activeOrders, id)
	}
}

func (s *Strategy) update(orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) {
	s.cancelOrders(session)

	switch s.Side {
	case "buy":
		s.updateOrders(orderExecutor, types.SideTypeBuy)
	case "sell":
		s.updateOrders(orderExecutor, types.SideTypeSell)
	case "both":
		s.updateOrders(orderExecutor, types.SideTypeBuy)
		s.updateOrders(orderExecutor, types.SideTypeSell)

	default:
		log.Panicf("undefined side: %s", s.Side)
	}

	s.book.C.Drain(1*time.Second, 3*time.Second)
}

func (s *Strategy) updateOrders(orderExecutor bbgo.OrderExecutor, side types.SideType) {
	book := s.book.Copy()

	var pvs types.PriceVolumeSlice

	switch side {
	case types.SideTypeBuy:
		pvs = book.Bids
	case types.SideTypeSell:
		pvs = book.Asks
	}

	if pvs == nil || len(pvs) == 0 {
		log.Warn("empty bids or asks")
		return
	}

	log.Infof("placing order behind volume: %f", s.BehindVolume.Float64())

	idx := pvs.IndexByVolumeDepth(s.BehindVolume)
	if idx == -1 {
		// do not place orders
		log.Warn("depth is not enough")
		return
	}

	var depthPrice = pvs[idx].Price
	var orders = s.generateOrders(s.Symbol, side, depthPrice, s.PriceTick, s.BaseQuantity, s.NumOrders)
	if len(orders) == 0 {
		log.Warn("empty orders")
		return
	}

	createdOrders, err := orderExecutor.SubmitOrders(context.Background(), orders...)
	if err != nil {
		log.WithError(err).Errorf("order submit error")
		return
	}

	// add created orders to the list
	for i, o := range createdOrders {
		s.activeOrders[o.ClientOrderID] = createdOrders[i]
	}
}

func (s *Strategy) generateOrders(symbol string, side types.SideType, price, priceTick, baseVolume fixedpoint.Value, numOrders int) (orders []types.SubmitOrder) {
	var expBase = fixedpoint.NewFromFloat(0.0)

	switch side {
	case types.SideTypeBuy:
		if priceTick > 0 {
			priceTick = -priceTick
		}

	case types.SideTypeSell:
		if priceTick < 0 {
			priceTick = -priceTick
		}
	}

	for i := 0; i < numOrders; i++ {
		volume := math.Exp(expBase.Float64()) * baseVolume.Float64()

		// skip order less than 10usd
		if volume*price.Float64() < 10.0 {
			log.Warnf("amount too small (< 10usd). price=%f volume=%f amount=%f", price.Float64(), volume, volume*price.Float64())
			continue
		}

		orders = append(orders, types.SubmitOrder{
			Symbol:   symbol,
			Side:     side,
			Type:     types.OrderTypeLimit,
			Price:    price.Float64(),
			Quantity: volume,
		})

		log.Infof("%s order: %.2f @ %f", side, volume, price.Float64())

		if len(orders) >= numOrders {
			break
		}

		price = price + priceTick
		declog := math.Log10(math.Abs(priceTick.Float64()))
		expBase += fixedpoint.NewFromFloat(math.Pow10(-int(declog)) * math.Abs(priceTick.Float64()))
		// log.Infof("expBase: %f", expBase.Float64())
	}

	return orders
}
