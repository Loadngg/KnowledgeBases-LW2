package data

import "math"

type TVarData struct {
	Scope
}

func NewTVarData() *TVarData {
	return &TVarData{
		Scope{
			Start: 110.0,
			End:   150.0,
		},
	}
}

func (d *TVarData) Small(T float64) float64 {
	return math.Pow(T-200, 2) / 10000
}

func (d *TVarData) Medium(T float64) float64 {
	return 1 - math.Pow(125-T, 2)/1000
}

func (d *TVarData) Big(T float64) float64 {
	return math.Pow(T-50, 2) / 10000
}

func (d *TVarData) GetScope() Scope {
	return d.Scope
}
