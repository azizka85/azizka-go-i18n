package i18n

import (
	"strconv"
	"strings"

	"github.com/azizka85/azizka-go-i18n/options"
)

type Translator struct {
	data          *options.DataOptions
	globalContext map[string]string

	extension func(text string, num interface{}, formatting map[string]string, data map[string]interface{}) string
}

func (translator *Translator) Add(data *options.DataOptions) {
	if translator.data == nil {
		translator.data = data
	} else {
		for key, value := range data.Values {
			translator.data.Values[key] = value
		}

		translator.data.Contexts = append(translator.data.Contexts, data.Contexts...)
	}
}

func (translator *Translator) SetContext(key string, value string) {
	translator.globalContext[key] = value
}

func (translator *Translator) ClearContext(key string) {
	delete(translator.globalContext, key)
}

func (translator *Translator) Extend(
	extension func(text string, num interface{}, formatting map[string]string, data map[string]interface{}) string,
) {
	translator.extension = extension
}

func (translator *Translator) ResetData() {
	translator.data = nil
}

func (translator *Translator) ResetContext() {
	translator.globalContext = make(map[string]string)
}

func (translator *Translator) Reset() {
	translator.ResetData()
	translator.ResetContext()
}

func (translator *Translator) Translate(
	text string,
	input ...interface{},
) string {
	var num interface{} = nil
	var formatting map[string]string = nil
	var defaultNumOrFormatting interface{}
	var numOrFormattingOrContext interface{}
	var formattingOrContext interface{}

	argsLen := len(input)

	if argsLen > 0 {
		defaultNumOrFormatting = input[0]
	}

	if argsLen > 1 {
		numOrFormattingOrContext = input[1]
	}

	if argsLen > 2 {
		formattingOrContext = input[2]
	}

	context := translator.globalContext

	if data, ok := defaultNumOrFormatting.(map[string]string); ok {
		formatting = data

		if data, ok := numOrFormattingOrContext.(map[string]string); ok {
			context = data
		}
	} else if data, ok := defaultNumOrFormatting.(int); ok {
		num = data

		if data, ok := numOrFormattingOrContext.(map[string]string); ok {
			formatting = data
		}

		if data, ok := formattingOrContext.(map[string]string); ok {
			context = data
		}
	} else {
		if data, ok := numOrFormattingOrContext.(int); ok {
			num = data

			if data, ok := formattingOrContext.(map[string]string); ok {
				formatting = data
			}
		} else {
			if data, ok := numOrFormattingOrContext.(map[string]string); ok {
				formatting = data
			}

			if data, ok := formattingOrContext.(map[string]string); ok {
				context = data
			}
		}
	}

	return translator.TranslateText(text, num, formatting, context)
}

func (translator *Translator) TranslateText(
	text string,
	num interface{},
	formatting map[string]string,
	context map[string]string,
) string {
	if context == nil {
		context = translator.globalContext
	}

	if translator.data == nil {
		return UseOriginalText(text, num, formatting)
	}

	var contextData *options.ContextOptions = GetContextData(translator.data, context)

	var result string
	var ok bool = false

	if contextData != nil {
		result, ok = translator.FindTranslation(text, num, formatting, contextData.Values)
	}

	if !ok {
		result, ok = translator.FindTranslation(text, num, formatting, translator.data.Values)
	}

	if !ok {
		result = UseOriginalText(text, num, formatting)
	}

	return result
}

func (translator *Translator) FindTranslation(
	text string,
	num interface{},
	formatting map[string]string,
	data map[string]interface{},
) (string, bool) {
	if data != nil {
		if value, ok := data[text]; ok {
			if val, ok := value.(map[string]interface{}); ok {
				if translator.extension != nil {
					result := translator.extension(text, num, formatting, val)
					result = ApplyNumbers(result, num)

					return ApplyFormatting(result, formatting), true
				} else {
					return UseOriginalText(text, num, formatting), true
				}
			}

			if num == nil {
				if val, ok := value.(string); ok {
					return ApplyFormatting(val, formatting), true
				}
			} else if val, ok := value.([][3]interface{}); ok {
				for _, triple := range val {
					numVal, hasNum := num.(int)

					low, hasLow := triple[0].(int)
					high, hasHigh := triple[1].(int)
					text, _ := triple[2].(string)

					if !hasNum && !hasLow && !hasHigh ||
						hasNum &&
							(hasLow && numVal >= low && (!hasHigh || numVal <= high) ||
								!hasLow && hasHigh && numVal <= high) {
						result := ApplyNumbers(text, num)

						return ApplyFormatting(result, formatting), true
					}
				}
			}
		}
	}

	return "", false
}

func ApplyNumbers(str string, num interface{}) string {
	numVal, _ := num.(int)

	str = strings.ReplaceAll(str, "-%n", strconv.Itoa(-numVal))
	str = strings.ReplaceAll(str, "%n", strconv.Itoa(numVal))

	return str
}

func ApplyFormatting(text string, formatting map[string]string) string {
	for key, value := range formatting {
		tpl := "%{" + key + "}"
		text = strings.ReplaceAll(text, tpl, value)
	}

	return text
}

func GetContextData(data *options.DataOptions, context map[string]string) *options.ContextOptions {
	if data.Contexts == nil {
		return nil
	}

	for _, ctx := range data.Contexts {
		equal := true

		for key, value := range ctx.Matches {
			equal = equal && value == context[key]

			if !equal {
				break
			}
		}

		if equal {
			return &ctx
		}
	}

	return nil
}

func UseOriginalText(text string, num interface{}, formatting map[string]string) string {
	if num == nil {
		return ApplyFormatting(text, formatting)
	}

	numVal, _ := num.(int)

	return ApplyFormatting(strings.ReplaceAll(text, "%n", strconv.Itoa(numVal)), formatting)
}

func CreateTranslator(data *options.DataOptions) *Translator {
	var translator = &Translator{}

	translator.ResetContext()
	translator.Add(data)

	return translator
}
