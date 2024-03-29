// Code generated by "callbackgen -type Strategy"; DO NOT EDIT.

package dca2

import (
	"github.com/c9s/bbgo/pkg/types"
)

func (s *Strategy) OnPosition(cb func(*types.Position)) {
	s.positionCallbacks = append(s.positionCallbacks, cb)
}

func (s *Strategy) EmitPosition(position *types.Position) {
	for _, cb := range s.positionCallbacks {
		cb(position)
	}
}

func (s *Strategy) OnProfit(cb func(*ProfitStats)) {
	s.profitCallbacks = append(s.profitCallbacks, cb)
}

func (s *Strategy) EmitProfit(profitStats *ProfitStats) {
	for _, cb := range s.profitCallbacks {
		cb(profitStats)
	}
}
