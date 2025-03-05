package repository

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"lr2/internal/constants"
	"lr2/internal/utils"
)

type Repository struct {
	RulesPath string
}

type Condition struct {
	Variable string
	Value    string
}

type Rule struct {
	Apparatus  string
	Conditions [][]Condition
}

func New(rules_path string) *Repository {
	return &Repository{
		RulesPath: rules_path,
	}
}

func (r *Repository) GetRules() ([]Rule, error) {
	ruleStrings, err := utils.ReadFileLines(r.RulesPath)
	if err != nil {
		return nil, err
	}
	rules, err := r.parseRules(ruleStrings)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (r *Repository) parseRules(ruleStrings []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range ruleStrings {
		rule, err := r.parseRule(s)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

func (r *Repository) parseRule(ruleStr string) (*Rule, error) {
	parts := strings.SplitN(ruleStr, " ТО ", 2)
	if len(parts) != 2 {
		return nil, errors.New(constants.InvalidRuleFormatError.String())
	}

	conditionPart := strings.TrimSpace(strings.TrimPrefix(parts[0], "ЕСЛИ "))
	apparatusPart := strings.TrimSpace(parts[1])

	conditionPart = strings.ReplaceAll(conditionPart, " И ", "И")
	conditionPart = strings.ReplaceAll(conditionPart, " ИЛИ ", "ИЛИ")

	conditions, err := r.parseConditions(conditionPart)
	if err != nil {
		return nil, err
	}
	return &Rule{
		Apparatus:  apparatusPart,
		Conditions: conditions,
	}, nil
}

func (r *Repository) parseConditions(s string) ([][]Condition, error) {
	var groups [][]Condition
	var currentGroup []Condition
	var buffer strings.Builder
	parenDepth := 0

	for i := 0; i < len(s); {
		rule, width := utf8.DecodeRuneInString(s[i:])
		i += width

		switch {
		case rule == '(':
			parenDepth++
			if parenDepth == 1 {
				if buffer.Len() > 0 {
					err := r.processBuffer(&buffer, &currentGroup)
					if err != nil {
						panic(err)
					}
				}
				continue
			}
		case rule == ')':
			parenDepth--
			if parenDepth == 0 {
				subConditions, err := r.parseConditions(buffer.String())
				if err != nil {
					return nil, err
				}
				for _, group := range subConditions {
					currentGroup = append(currentGroup, group...)
				}
				buffer.Reset()
				continue
			}
		case strings.HasPrefix(s[i-width:], "ИЛИ") && parenDepth == 0:
			err := r.processBuffer(&buffer, &currentGroup)
			if err != nil {
				return nil, err
			}
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = nil
			}
			i += 2
			continue
		case rule == 'И' && parenDepth == 0:
			err := r.processBuffer(&buffer, &currentGroup)
			if err != nil {
				return nil, err
			}
			continue
		}

		buffer.WriteRune(rule)
	}

	err := r.processBuffer(&buffer, &currentGroup)
	if err != nil {
		return nil, err
	}
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups, nil
}

func (r *Repository) processBuffer(buffer *strings.Builder, group *[]Condition) error {
	if buffer.Len() == 0 {
		return nil
	}

	condStr := strings.TrimFunc(buffer.String(), func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')'
	})

	if condStr != "" {
		cond, err := r.parseSingleCondition(condStr)
		if err != nil {
			return err
		}
		*group = append(*group, cond)
	}
	buffer.Reset()
	return nil
}

func (r *Repository) parseSingleCondition(condStr string) (Condition, error) {
	parts := strings.SplitN(condStr, "=", 2)
	if len(parts) != 2 {
		return Condition{}, fmt.Errorf(constants.InvalidConditionError.String(), condStr)
	}

	var variable, value string
	switch strings.TrimSpace(parts[0]) {
	case "расход_сырья":
		variable = "G"
	case "температура_процесса":
		variable = "T"
	default:
		return Condition{}, fmt.Errorf(constants.UnknownVariableError.String(), parts[0])
	}

	value = strings.TrimSpace(parts[1])
	switch value {
	case "слегка_малый", "слегка_малая":
		value = "SlightlySmall"
	case "малый", "малая":
		value = "Small"
	case "средний", "средняя":
		value = "Medium"
	case "большой", "большая":
		value = "Big"
	default:
		return Condition{}, fmt.Errorf(constants.UnknownValueError.String(), value)
	}

	return Condition{
		Variable: variable,
		Value:    value,
	}, nil
}
