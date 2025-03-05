package constants

type TextEnum int

const (
	WindowName TextEnum = iota
	GEntryLabel
	TEntryLabel
	ShowCharts
	ChartsLabel
	Small
	Medium
	Big
	StartServerError
	ValueEmptyError
	IncorrectValueError
	RuleWeightLabel
	Result
	ValueNotInScope
	InvalidRuleFormatError
	InvalidConditionError
	UnknownVariableError
	UnknownValueError
)

var textLabels = map[TextEnum]string{
	WindowName:             "Базы знаний Лр2",
	GEntryLabel:            "Расход сырья",
	TEntryLabel:            "t℃ процесса",
	ShowCharts:             "Показать графики",
	ChartsLabel:            "Графики расхода сырья и t℃ процесса",
	Small:                  "Малый",
	Medium:                 "Средний",
	Big:                    "Большой",
	StartServerError:       "Ошибка запуска сервера: %v",
	ValueEmptyError:        "Значение не может быть пустым",
	IncorrectValueError:    "Значение должно быть числом",
	RuleWeightLabel:        "Правило %d: вес %.2f\n",
	Result:                 "\nВыбранные аппараты: %s\n",
	ValueNotInScope:        "Значение %s не входит в диапазон %d - %d",
	InvalidRuleFormatError: "Неверный формат правила",
	InvalidConditionError:  "Неверное условие: %s",
	UnknownVariableError:   "Неизвестная переменная: %s",
	UnknownValueError:      "Неизвестное значение: %s",
}

func (e TextEnum) String() string {
	if val, ok := textLabels[e]; ok {
		return val
	}
	return "Неизвестный ключ"
}
