package data

type Scope struct {
	Start float64
	End   float64
}

type Var interface {
	Small(float64) float64
	Medium(float64) float64
	Big(float64) float64
	GetScope() Scope
}

type Data struct {
	G Var
	T Var
}

func New() *Data {
	return &Data{
		G: NewGVarData(),
		T: NewTVarData(),
	}
}
