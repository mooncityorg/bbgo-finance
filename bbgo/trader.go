package bbgo

import (
	"context"
	"fmt"
	"github.com/c9s/bbgo/pkg/slack/slackstyle"
	"github.com/c9s/bbgo/pkg/util"
	"github.com/leekchan/accounting"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"strconv"
	"time"
)

var USD = accounting.Accounting{Symbol: "$ ", Precision: 2}
var BTC = accounting.Accounting{Symbol: "BTC ", Precision: 8}

type Trader struct {
	// Context is trading Context
	Context *TradingContext

	Exchange *BinanceExchange

	Slack *slack.Client

	TradingChannel string
	ErrorChannel   string
	InfoChannel    string
}

func (t *Trader) Infof(format string, args ...interface{}) {
	var slackAttachments []slack.Attachment = nil
	var slackArgsStartIdx = -1
	for idx, arg := range args {
		switch a := arg.(type) {

		// concrete type assert first
		case slack.Attachment:
			if slackArgsStartIdx == -1 {
				slackArgsStartIdx = idx
			}
			slackAttachments = append(slackAttachments, a)

		case slackstyle.SlackAttachmentCreator:
			if slackArgsStartIdx == -1 {
				slackArgsStartIdx = idx
			}
			slackAttachments = append(slackAttachments, a.SlackAttachment())

		}
	}

	var nonSlackArgs = []interface{}{}
	if slackArgsStartIdx > 0 {
		nonSlackArgs = args[:slackArgsStartIdx]
	}

	logrus.Infof(format, nonSlackArgs...)

	_, _, err := t.Slack.PostMessageContext(context.Background(), t.InfoChannel,
		slack.MsgOptionText(fmt.Sprintf(format, nonSlackArgs...), true),
		slack.MsgOptionAttachments(slackAttachments...))
	if err != nil {
		logrus.WithError(err).Error("Slack error:", err)
	}
}

func (t *Trader) Errorf(err error, format string, args ...interface{}) {
	logrus.WithError(err).Errorf(format, args...)
	_, _, err2 := t.Slack.PostMessageContext(context.Background(), t.ErrorChannel,
		slack.MsgOptionText("ERROR: "+err.Error()+" "+fmt.Sprintf(format, args...), true))
	if err2 != nil {
		logrus.WithError(err2).Error("Slack error:", err2)
	}
}

func (t *Trader) ReportTrade(e *BinanceExecutionReportEvent, trade *Trade) {
	var color = ""
	if trade.IsBuyer {
		color = "#228B22"
	} else {
		color = "#DC143C"
	}

	_, _, err := t.Slack.PostMessageContext(context.Background(), t.TradingChannel,
		slack.MsgOptionText(util.Render(`:handshake: {{ .CurrentExecutionType }} execution`, e), true),
		slack.MsgOptionAttachments(slack.Attachment{
			Title: "New Trade",
			Color: color,
			// Pretext:       "",
			// Text:          "",
			Fields: []slack.AttachmentField{
				{Title: "Market", Value: trade.Market, Short: true,},
				{Title: "Side", Value: e.Side, Short: true,},
				{Title: "Price", Value: USD.FormatMoney(trade.Price), Short: true,},
				{Title: "Volume", Value: t.Context.Market.FormatVolume(trade.Volume), Short: true,},
			},
			// Footer:     tradingCtx.TradeStartTime.Format(time.RFC822),
			// FooterIcon: "",
		}))

	if err != nil {
		t.Errorf(err, "Slack send error")
	}
}

func (t *Trader) ReportPnL() {
	t.Context.UpdatePnL()

	tradingCtx := t.Context
	logrus.Infof("current price:  %s", USD.FormatMoneyFloat64(tradingCtx.CurrentPrice))
	logrus.Infof("average bid price:  %s", USD.FormatMoneyFloat64(tradingCtx.AverageBidPrice))
	logrus.Infof("Stock volume:   %f", tradingCtx.Stock)
	logrus.Infof("overall Profit:  %s", USD.FormatMoneyFloat64(tradingCtx.Profit))

	var color = ""
	if tradingCtx.Profit > 0 {
		color = slackstyle.Green
	} else {
		color = slackstyle.Red
	}

	_, _, err := t.Slack.PostMessageContext(context.Background(), t.TradingChannel,
		slack.MsgOptionText(util.Render(
			`:heavy_dollar_sign: Here is your *{{ .Symbol }}* PnL report collected since *{{ .startTime }}*`,
			map[string]interface{}{
				"Symbol":    tradingCtx.Symbol,
				"startTime": tradingCtx.TradeStartTime.Format(time.RFC822),
			}), true),
		slack.MsgOptionAttachments(slack.Attachment{
			Title: "Profit and Loss report",
			Color: color,
			// Pretext:       "",
			// Text:          "",
			Fields: []slack.AttachmentField{
				{
					Title: "Market",
					Value: tradingCtx.Symbol,
					Short: true,
				},
				{
					Title: "Profit",
					Value: USD.FormatMoney(tradingCtx.Profit),
					Short: true,
				},
				{
					Title: "Current Price",
					Value: USD.FormatMoney(tradingCtx.CurrentPrice),
					Short: true,
				},
				{
					Title: "Average Bid Price",
					Value: USD.FormatMoney(tradingCtx.AverageBidPrice),
					Short: true,
				},
				{
					Title: "Current Stock",
					Value: tradingCtx.Market.FormatVolume(tradingCtx.Stock),
					Short: true,
				},
				{
					Title: "Number of Trades",
					Value: strconv.Itoa(len(tradingCtx.Trades)),
					Short: true,
				},
			},
			Footer:     tradingCtx.TradeStartTime.Format(time.RFC822),
			FooterIcon: "",
		}))

	if err != nil {
		t.Errorf(err, "Slack send error")
	}
}

func (t *Trader) SubmitOrder(ctx context.Context, order *Order) {
	t.Infof(":memo: Submitting %s order on side %s with volume: %s", order.Type, order.Side, order.VolumeStr, order.SlackAttachment())

	err := t.Exchange.SubmitOrder(ctx, order)
	if err != nil {
		t.Errorf(err, "order create error: side %s volume: %s", order.Side, order.VolumeStr)
		return
	}
}

