package kucoin

import (
	"context"
	"strconv"
	"time"

	"github.com/c9s/bbgo/pkg/exchange/kucoin/kucoinapi"
	"github.com/c9s/bbgo/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ErrMissingSequence = errors.New("sequence is missing")

// OKB is the platform currency of OKEx, pre-allocate static string here
const KCS = "KCS"

var log = logrus.WithFields(logrus.Fields{
	"exchange": "kucoin",
})

type Exchange struct {
	key, secret, passphrase string
	client                  *kucoinapi.RestClient
}

func New(key, secret, passphrase string) *Exchange {
	client := kucoinapi.NewClient()

	// for public access mode
	if len(key) > 0 && len(secret) > 0 && len(passphrase) > 0 {
		client.Auth(key, secret, passphrase)
	}

	return &Exchange{
		key:        key,
		secret:     secret,
		passphrase: passphrase,
		client:     client,
	}
}

func (e *Exchange) Name() types.ExchangeName {
	return types.ExchangeKucoin
}

func (e *Exchange) PlatformFeeCurrency() string {
	return KCS
}

func (e *Exchange) QueryAccount(ctx context.Context) (*types.Account, error) {
	accounts, err := e.client.AccountService.ListAccounts()
	if err != nil {
		return nil, err
	}

	// for now, we only return the trading account
	a := types.NewAccount()
	balances := toGlobalBalanceMap(accounts)
	a.UpdateBalances(balances)
	return a, nil
}

func (e *Exchange) QueryAccountBalances(ctx context.Context) (types.BalanceMap, error) {
	accounts, err := e.client.AccountService.ListAccounts()
	if err != nil {
		return nil, err
	}

	return toGlobalBalanceMap(accounts), nil
}

func (e *Exchange) QueryMarkets(ctx context.Context) (types.MarketMap, error) {
	markets, err := e.client.MarketDataService.ListSymbols()
	if err != nil {
		return nil, err
	}

	marketMap := types.MarketMap{}
	for _, s := range markets {
		market := toGlobalMarket(s)
		marketMap.Add(market)
	}

	return marketMap, nil
}

func (e *Exchange) QueryTicker(ctx context.Context, symbol string) (*types.Ticker, error) {
	s, err := e.client.MarketDataService.GetTicker24HStat(symbol)
	if err != nil {
		return nil, err
	}

	ticker := toGlobalTicker(*s)
	return &ticker, nil
}

func (e *Exchange) QueryTickers(ctx context.Context, symbols ...string) (map[string]types.Ticker, error) {
	tickers := map[string]types.Ticker{}
	if len(symbols) > 0 {
		for _, s := range symbols {
			t, err := e.QueryTicker(ctx, s)
			if err != nil {
				return nil, err
			}

			tickers[s] = *t
		}

		return tickers, nil
	}

	allTickers, err := e.client.MarketDataService.ListTickers()
	if err != nil {
		return nil, err
	}

	for _, s := range allTickers.Ticker {
		tickers[s.Symbol] = toGlobalTicker(s)
	}

	return tickers, nil
}

func (e *Exchange) QueryKLines(ctx context.Context, symbol string, interval types.Interval, options types.KLineQueryOptions) ([]types.KLine, error) {
	panic("implement me")
}

func (e *Exchange) SubmitOrders(ctx context.Context, orders ...types.SubmitOrder) (createdOrders types.OrderSlice, err error) {
	for _, order := range orders {
		req := e.client.TradeService.NewPlaceOrderRequest()
		req.Symbol(toLocalSymbol(order.Symbol))
		req.Side(toLocalSide(order.Side))

		if order.ClientOrderID != "" {
			req.ClientOrderID(order.ClientOrderID)
		}

		if len(order.QuantityString) > 0 {
			req.Size(order.QuantityString)
		} else if order.Market.Symbol != "" {
			req.Size(order.Market.FormatQuantity(order.Quantity))
		} else {
			req.Size(strconv.FormatFloat(order.Quantity, 'f', 8, 64))
		}

		// set price field for limit orders
		switch order.Type {
		case types.OrderTypeStopLimit, types.OrderTypeLimit:
			if len(order.PriceString) > 0 {
				req.Price(order.PriceString)
			} else if order.Market.Symbol != "" {
				req.Price(order.Market.FormatPrice(order.Price))
			}
		}

		switch order.TimeInForce {
		case "FOK":
			req.TimeInForce(kucoinapi.TimeInForceFOK)
		case "IOC":
			req.TimeInForce(kucoinapi.TimeInForceIOC)
		default:
			// default to GTC
			req.TimeInForce(kucoinapi.TimeInForceGTC)
		}

		orderResponse, err := req.Do(ctx)
		if err != nil {
			return createdOrders, err
		}

		createdOrders = append(createdOrders, types.Order{
			SubmitOrder:      order,
			Exchange:         types.ExchangeKucoin,
			OrderID:          hashStringID(orderResponse.OrderID),
			Status:           types.OrderStatusNew,
			ExecutedQuantity: 0,
			IsWorking:        true,
			CreationTime:     types.Time(time.Now()),
			UpdateTime:       types.Time(time.Now()),
		})
	}

	return createdOrders, err
}

func (e *Exchange) QueryOpenOrders(ctx context.Context, symbol string) (orders []types.Order, err error) {
	req := e.client.TradeService.NewListOrdersRequest()
	req.Symbol(toLocalSymbol(symbol))

	return nil, nil
}

func (e *Exchange) CancelOrders(ctx context.Context, orders ...types.Order) error {
	panic("implement me")
	return nil
}

func (e *Exchange) NewStream() types.Stream {
	return NewStream(e.client, e)
}

func (e *Exchange) QueryDepth(ctx context.Context, symbol string) (types.SliceOrderBook, int64, error) {
	orderBook, err := e.client.MarketDataService.GetOrderBook(toLocalSymbol(symbol), 100)
	if err != nil {
		return types.SliceOrderBook{}, 0, err
	}

	if len(orderBook.Sequence) == 0 {
		return types.SliceOrderBook{}, 0, ErrMissingSequence
	}

	sequence, err := strconv.ParseInt(orderBook.Sequence, 10, 64)
	if err != nil {
		return types.SliceOrderBook{}, 0, err
	}

	return types.SliceOrderBook{
		Symbol: toGlobalSymbol(symbol),
		Bids:   orderBook.Bids,
		Asks:   orderBook.Asks,
	}, sequence, nil
}
