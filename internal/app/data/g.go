package data

import (
	"math"

	"lr2/internal/utils"
)

type GVarData struct {
	Scope
}

func NewGVarData() *GVarData {
	return &GVarData{
		Scope: Scope{
			Start: 70.0,
			End:   110.0,
		},
	}
}

func (d *GVarData) Small(G float64) float64 {
	return utils.Clamp(math.Pow(math.Exp(-0.2*math.Log(10*math.Abs(G-75.1))), 2))
}

func (d *GVarData) Medium(G float64) float64 {
	return utils.Clamp(math.Pow(math.Exp(-0.2*math.Log(10*math.Abs(G-85.1))), 2))
}

func (d *GVarData) Big(G float64) float64 {
	return utils.Clamp(math.Pow(math.Exp(-0.2*math.Log(10*math.Abs(G-100.1))), 2))
}

func (d *GVarData) GetScope() Scope {
	return d.Scope
}
