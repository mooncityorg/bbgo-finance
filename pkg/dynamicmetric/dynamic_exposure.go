package dynamicmetric

import (
	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/indicator"
	"github.com/c9s/bbgo/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math"
)

type DynamicExposure struct {
	// BollBandExposure calculates the max exposure with the Bollinger Band
	BollBandExposure *DynamicExposureBollBand `json:"bollBandExposure"`
}

// Initialize dynamic exposure
func (d *DynamicExposure) Initialize(symbol string, session *bbgo.ExchangeSession, standardIndicatorSet *bbgo.StandardIndicatorSet) {
	switch {
	case d.BollBandExposure != nil:
		d.BollBandExposure.initialize(symbol, session, standardIndicatorSet)
	}
}

func (d *DynamicExposure) IsEnabled() bool {
	return d.BollBandExposure != nil
}

// GetMaxExposure returns the max exposure
func (d *DynamicExposure) GetMaxExposure(price float64) (maxExposure fixedpoint.Value, err error) {
	switch {
	case d.BollBandExposure != nil:
		return d.BollBandExposure.getMaxExposure(price)
	default:
		return fixedpoint.Zero, errors.New("dynamic exposure is not enabled")
	}
}

// DynamicExposureBollBand calculates the max exposure with the Bollinger Band
type DynamicExposureBollBand struct {
	// DynamicExposureBollBandScale is used to define the exposure range with the given percentage.
	DynamicExposureBollBandScale *bbgo.PercentageScale `json:"dynamicExposurePositionScale"`

	types.IntervalWindowBandWidth

	StandardIndicatorSet *bbgo.StandardIndicatorSet

	dynamicExposureBollBand *indicator.BOLL
}

// initialize dynamic exposure with Bollinger Band
func (d *DynamicExposureBollBand) initialize(symbol string, session *bbgo.ExchangeSession, standardIndicatorSet *bbgo.StandardIndicatorSet) {
	d.StandardIndicatorSet = standardIndicatorSet
	d.dynamicExposureBollBand = d.StandardIndicatorSet.BOLL(d.IntervalWindow, d.BandWidth)

	// Subscribe kline
	session.Subscribe(types.KLineChannel, symbol, types.SubscribeOptions{
		Interval: d.dynamicExposureBollBand.Interval,
	})
}

// getMaxExposure returns the max exposure
func (d *DynamicExposureBollBand) getMaxExposure(price float64) (fixedpoint.Value, error) {
	downBand := d.dynamicExposureBollBand.DownBand.Last()
	upBand := d.dynamicExposureBollBand.UpBand.Last()
	sma := d.dynamicExposureBollBand.SMA.Last()
	log.Infof("dynamicExposureBollBand bollinger band: up %f sma %f down %f", upBand, sma, downBand)

	bandPercentage := 0.0
	if price < sma {
		// should be negative percentage
		bandPercentage = (price - sma) / math.Abs(sma-downBand)
	} else if price > sma {
		// should be positive percentage
		bandPercentage = (price - sma) / math.Abs(upBand-sma)
	}

	v, err := d.DynamicExposureBollBandScale.Scale(bandPercentage)
	if err != nil {
		return fixedpoint.Zero, err
	}
	return fixedpoint.NewFromFloat(v), nil
}
