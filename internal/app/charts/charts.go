package charts

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"lr2/internal/app/data"
	"lr2/internal/constants"
)

type Charts struct {
	ChartFile string
	Data      *data.Data
}

type result struct {
	Small  []opts.LineData
	Medium []opts.LineData
	Big    []opts.LineData
}

func New(chartFile string, d *data.Data) *Charts {
	return &Charts{
		ChartFile: chartFile,
		Data:      d,
	}
}

func (c *Charts) Generate() {
	page := components.NewPage()
	page.SetPageTitle(constants.ChartsLabel.String())

	gLine := generateChart(constants.GEntryLabel.String(), c.Data.G)
	tLine := generateChart(constants.TEntryLabel.String(), c.Data.T)

	page.AddCharts(gLine, tLine)
	f, err := os.Create(c.ChartFile)
	if err != nil {
		panic(err)

	}
	page.Render(f)
}

func generateChart(title string, v data.Var) *charts.Line {
	line := charts.NewLine()

	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	result := generateValues(v)
	line.SetXAxis(generateScopeAxis(v)).
		AddSeries(constants.Small.String(), result.Small).
		AddSeries(constants.Medium.String(), result.Medium).
		AddSeries(constants.Big.String(), result.Big).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	return line
}

func generateScopeAxis(d data.Var) []int {
	scope := d.GetScope()

	var axis []int
	for T := scope.Start; T <= scope.End; T += 0.1 {
		axis = append(axis, int(T))
	}
	return axis
}

func generateValues(d data.Var) *result {
	scope := d.GetScope()

	var result result
	for i := scope.Start; i <= scope.End; i += 0.1 {
		result.Small = append(result.Small, opts.LineData{Value: d.Small(i)})
		result.Medium = append(result.Medium, opts.LineData{Value: d.Medium(i)})
		result.Big = append(result.Big, opts.LineData{Value: d.Big(i)})
	}
	return &result
}
