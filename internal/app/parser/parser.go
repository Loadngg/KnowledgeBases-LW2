package parser

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"lr2/internal/app/data"
	"lr2/internal/app/repository"
	"lr2/internal/constants"
	"lr2/internal/utils"
)

type Parser struct {
	Repository *repository.Repository
	Data       *data.Data
}

type MuValues struct {
	G             []float64
	T             []float64
	RulesMuMatrix [][]int
}

type MuParsableValues struct {
	GSlightlySmall float64
	GSmall         float64
	GMedium        float64
	GBig           float64
	TSmall         float64
	TMedium        float64
	TBig           float64
}

type ApparatusWeight struct {
	Name   string
	Weight float64
}

func New(r *repository.Repository, d *data.Data) *Parser {
	return &Parser{
		Repository: r,
		Data:       d,
	}
}

func (p *Parser) Parse(g float64, t float64) (*string, error) {
	muValues, err := p.generateData(g, t)
	if err != nil {
		return nil, err
	}

	muParsableValues := &MuParsableValues{
		GSlightlySmall: math.Sqrt(muValues.G[0]),
		GSmall:         muValues.G[0],
		GMedium:        muValues.G[1],
		GBig:           muValues.G[2],
		TSmall:         muValues.T[0],
		TMedium:        muValues.T[1],
		TBig:           muValues.T[2],
	}

	rules, err := p.Repository.GetRules()
	if err != nil {
		return nil, err
	}

	var output strings.Builder
	ruleWeights := make([]float64, len(rules))
	for i, rule := range rules {
		weight := p.calculateRuleWeight(rule, *muParsableValues)
		ruleWeights[i] = utils.RoundValue(weight)
		output.WriteString(fmt.Sprintf(constants.RuleWeightLabel.String(), i+1, ruleWeights[i]))
	}

	apparatusWeights := p.calculateApparatusWeights(ruleWeights, muValues.RulesMuMatrix)
	sort.Slice(apparatusWeights, func(i, j int) bool {
		return apparatusWeights[i].Weight > apparatusWeights[j].Weight
	})

	maxWeight := apparatusWeights[0].Weight
	var result []string
	for _, aw := range apparatusWeights {
		if aw.Weight == maxWeight {
			result = append(result, aw.Name)
		}
	}

	sort.Strings(result)
	output.WriteString(fmt.Sprintf(constants.Result.String(), strings.Join(result, ", ")))
	outputStr := output.String()
	return &outputStr, nil
}

func (p *Parser) checkScope(f float64, v data.Var) bool {
	scope := v.GetScope()
	return f >= scope.Start && f <= scope.End
}

func (p *Parser) generateData(g float64, t float64) (*MuValues, error) {
	if !p.checkScope(g, p.Data.G) {
		return nil, fmt.Errorf(constants.ValueNotInScope.String(), "G", p.Data.G.GetScope().Start, p.Data.G.GetScope().End)
	}
	if !p.checkScope(t, p.Data.T) {
		return nil, fmt.Errorf(constants.ValueNotInScope.String(), "T", p.Data.T.GetScope().Start, p.Data.T.GetScope().End)
	}

	gMuValues := []float64{
		utils.RoundValue(p.Data.G.Small(g)),
		utils.RoundValue(p.Data.G.Medium(g)),
		utils.RoundValue(p.Data.G.Big(g)),
	}

	tMuValues := []float64{
		utils.RoundValue(p.Data.T.Small(t)),
		utils.RoundValue(p.Data.T.Medium(t)),
		utils.RoundValue(p.Data.T.Big(t)),
	}

	rulesMuMatrix := [][]int{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	return &MuValues{
		G:             gMuValues,
		T:             tMuValues,
		RulesMuMatrix: rulesMuMatrix,
	}, nil
}

func (p *Parser) calculateRuleWeight(rule repository.Rule, mu MuParsableValues) float64 {
	maxWeight := 0.0
	for _, conditionGroup := range rule.Conditions {
		groupWeight := 1.0
		for _, cond := range conditionGroup {
			var val float64
			switch cond.Variable {
			case "G":
				switch cond.Value {
				case "SlightlySmall":
					val = mu.GSlightlySmall
				case "Small":
					val = mu.GSmall
				case "Medium":
					val = mu.GMedium
				case "Big":
					val = mu.GBig
				}
			case "T":
				switch cond.Value {
				case "Small":
					val = mu.TSmall
				case "Medium":
					val = mu.TMedium
				case "Big":
					val = mu.TBig
				}
			}
			groupWeight = math.Min(groupWeight, val)
		}
		maxWeight = math.Max(maxWeight, groupWeight)
	}
	return maxWeight
}

func (p *Parser) calculateApparatusWeights(ruleWeights []float64, matrix [][]int) []ApparatusWeight {
	apparatusNames := []string{"A", "B", "C", "D"}
	weights := make([]ApparatusWeight, len(apparatusNames))

	for apparatusIdx := range apparatusNames {
		minVal := 1.0

		for ruleIdx, w := range ruleWeights {
			term := 1 - w + float64(matrix[ruleIdx][apparatusIdx])
			minVal = math.Min(minVal, term)
		}

		weights[apparatusIdx] = ApparatusWeight{
			Name:   apparatusNames[apparatusIdx],
			Weight: minVal,
		}
	}

	return weights
}
