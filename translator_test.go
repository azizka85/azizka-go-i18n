package i18n

import (
	"azizka-go/i18n/options"
	"testing"
)

func TestTranslateHello(t *testing.T) {
	key := "Hello"
	value := "Hello translated"

	values := make(map[string]interface{})

	values[key] = value

	translator := CreateTranslator(&options.DataOptions{
		Values:   values,
		Contexts: nil,
	})

	actual := translator.Translate(key, nil, nil, nil)

	if actual != value {
		t.Errorf("Should translate '%v' to '%v' but the result is '%v'", key, value, actual)
	}
}

func TestTranslatePluralText(t *testing.T) {
	key := "%n comments"

	zeroComments := "0 comments"
	oneComment := "1 comment"
	twoComments := "2 comments"
	tenComments := "10 comments"

	values := make(map[string]interface{})

	values[key] = [][3]interface{}{
		{0, 0, "%n comments"},
		{1, 1, "%n comment"},
		{2, nil, "%n comments"},
	}

	translator := CreateTranslator(&options.DataOptions{
		Values:   values,
		Contexts: nil,
	})

	actual := translator.Translate(key, 0, nil, nil)

	if actual != zeroComments {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 0, zeroComments, actual)
	}

	actual = translator.Translate(key, 1, nil, nil)

	if actual != oneComment {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 1, oneComment, actual)
	}

	actual = translator.Translate(key, 2, nil, nil)

	if actual != twoComments {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 2, twoComments, actual)
	}

	actual = translator.Translate(key, 10, nil, nil)

	if actual != tenComments {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 10, tenComments, actual)
	}
}

func TestTranslatePluralTextWithNegativeNumber(t *testing.T) {
	key := "Due in %n days"

	dueTenDaysAgo := "Due 10 days ago"
	dueTwoDaysAgo := "Due 2 days ago"
	dueYesterday := "Due Yesterday"
	dueToday := "Due Today"
	dueTomorrow := "Due Tomorrow"
	dueInTwoDays := "Due in 2 days"
	dueInTenDays := "Due in 10 days"

	values := make(map[string]interface{})

	values[key] = [][3]interface{}{
		{nil, -2, "Due -%n days ago"},
		{-1, -1, "Due Yesterday"},
		{0, 0, "Due Today"},
		{1, 1, "Due Tomorrow"},
		{2, nil, "Due in %n days"},
	}

	translator := CreateTranslator(&options.DataOptions{
		Values:   values,
		Contexts: nil,
	})

	actual := translator.Translate(key, -10, nil, nil)

	if actual != dueTenDaysAgo {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, -10, dueTenDaysAgo, actual)
	}

	actual = translator.Translate(key, -2, nil, nil)

	if actual != dueTwoDaysAgo {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, -2, dueTwoDaysAgo, actual)
	}

	actual = translator.Translate(key, -1, nil, nil)

	if actual != dueYesterday {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, -1, dueYesterday, actual)
	}

	actual = translator.Translate(key, 0, nil, nil)

	if actual != dueToday {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 0, dueToday, actual)
	}

	actual = translator.Translate(key, 1, nil, nil)

	if actual != dueTomorrow {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 1, dueTomorrow, actual)
	}

	actual = translator.Translate(key, 2, nil, nil)

	if actual != dueInTwoDays {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 2, dueInTwoDays, actual)
	}

	actual = translator.Translate(key, 10, nil, nil)

	if actual != dueInTenDays {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 10, dueInTenDays, actual)
	}
}

func TestTranslateTextWithFormatting(t *testing.T) {
	key := "Welcome %{name}"
	value := "Welcome John"

	translator := CreateTranslator(nil)

	actual := translator.Translate(
		key,
		map[string]string{
			"name": "John",
		},
		nil,
		nil,
	)

	if actual != value {
		t.Errorf("Should translate '%v' with name '%v' to '%v' but the result is '%v'", key, "John", value, actual)
	}
}

func TestTranslateTextUsingContexts(t *testing.T) {
	key := "%{name} updated their profile"

	johnValue := "John updated his profile"
	janeValue := "Jane updated her profile"

	maleValues := make(map[string]interface{})

	maleValues[key] = "%{name} updated his profile"

	femaleValues := make(map[string]interface{})

	femaleValues[key] = "%{name} updated her profile"

	contexts := []options.ContextOptions{{
		Matches: map[string]string{
			"gender": "male",
		},
		Values: maleValues,
	}, {
		Matches: map[string]string{
			"gender": "female",
		},
		Values: femaleValues,
	}}

	translator := CreateTranslator(&options.DataOptions{
		Values:   nil,
		Contexts: contexts,
	})

	actual := translator.Translate(
		key,
		map[string]string{
			"name": "John",
		},
		map[string]string{
			"gender": "male",
		},
		nil,
	)

	if actual != johnValue {
		t.Errorf(
			"Should translate '%v' with name '%v' and gender '%v' to '%v' but the result is '%v'",
			key,
			"John",
			"male",
			johnValue,
			actual,
		)
	}

	actual = translator.Translate(
		key,
		map[string]string{
			"name": "Jane",
		},
		map[string]string{
			"gender": "female",
		},
		nil,
	)

	if actual != janeValue {
		t.Errorf(
			"Should translate '%v' with name '%v' and gender '%v' to '%v' but the result is '%v'",
			key,
			"Jane",
			"female",
			janeValue,
			actual,
		)
	}
}

func TestTranslatePluralTextUsingContexts(t *testing.T) {
	key := "%{name} uploaded %n photos to their %{album} album"

	johnValue := "John uploaded 1 photo to his Buck's Night album"
	janeValue := "Jane uploaded 4 photos to her Hen's Night album"

	maleValues := make(map[string]interface{})

	maleValues[key] = [][3]interface{}{
		{0, 0, "%{name} uploaded %n photos to his %{album} album"},
		{1, 1, "%{name} uploaded %n photo to his %{album} album"},
		{2, nil, "%{name} uploaded %n photos to his %{album} album"},
	}

	femaleValues := make(map[string]interface{})

	femaleValues[key] = [][3]interface{}{
		{0, 0, "%{name} uploaded %n photos to her %{album} album"},
		{1, 1, "%{name} uploaded %n photo to her %{album} album"},
		{2, nil, "%{name} uploaded %n photos to her %{album} album"},
	}

	contexts := []options.ContextOptions{{
		Matches: map[string]string{
			"gender": "male",
		},
		Values: maleValues,
	}, {
		Matches: map[string]string{
			"gender": "female",
		},
		Values: femaleValues,
	}}

	translator := CreateTranslator(&options.DataOptions{
		Values:   nil,
		Contexts: contexts,
	})

	actual := translator.Translate(
		key,
		1,
		map[string]string{
			"name":  "John",
			"album": "Buck's Night",
		},
		map[string]string{
			"gender": "male",
		},
	)

	if actual != johnValue {
		t.Errorf(
			"Should translate '%v' with name '%v', album '%v', num '%v' and gender '%v' to '%v' but the result is '%v'",
			key,
			"John",
			"Buck's Night",
			1,
			"male",
			johnValue,
			actual,
		)
	}

	actual = translator.Translate(
		key,
		4,
		map[string]string{
			"name":  "Jane",
			"album": "Hen's Night",
		},
		map[string]string{
			"gender": "female",
		},
	)

	if actual != janeValue {
		t.Errorf(
			"Should translate '%v' with name '%v', album '%v', num '%v' and gender '%v' to '%v' but the result is '%v'",
			key,
			"Jane",
			"Hen's Night",
			4,
			"female",
			janeValue,
			actual,
		)
	}
}

func TestTranslatePluralTextUsingExtension(t *testing.T) {
	key := "%n results"

	zeroResults := "нет результатов"
	oneResult := "1 результат"
	elevenResults := "11 результатов"
	fourResults := "4 результата"
	results := "101 результат"

	values := make(map[string]interface{})

	values[key] = map[string]interface{}{
		"zero":  "нет результатов",
		"one":   "%n результат",
		"few":   "%n результата",
		"many":  "%n результатов",
		"other": "%n результаты",
	}

	translator := CreateTranslator(&options.DataOptions{
		Values:   values,
		Contexts: nil,
	})

	getPluralisationKey := func(num interface{}) string {
		numVal, ok := num.(int)

		if !ok || numVal == 0 {
			return "zero"
		}

		if numVal%10 == 1 && numVal%100 != 11 {
			return "one"
		}

		if (numVal%10 == 2 || numVal%10 == 3 || numVal%10 == 4) &&
			numVal%100 != 12 && numVal%100 != 13 && numVal%100 != 14 {
			return "few"
		}

		if numVal%10 == 0 || numVal%10 == 5 || numVal%10 == 6 || numVal%10 == 7 || numVal%10 == 8 || numVal%10 == 9 ||
			numVal%100 == 11 || numVal%100 == 12 || numVal%100 == 13 || numVal%100 == 14 {
			return "many"
		}

		return "other"
	}

	russianExtension := func(
		text string,
		num interface{},
		formatting map[string]string,
		data map[string]interface{},
	) string {
		key := getPluralisationKey(num)

		return data[key].(string)
	}

	translator.Extend(russianExtension)

	actual := translator.Translate(key, 0, nil, nil)

	if actual != zeroResults {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 0, zeroResults, actual)
	}

	actual = translator.Translate(key, 1, nil, nil)

	if actual != oneResult {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 1, oneResult, actual)
	}

	actual = translator.Translate(key, 11, nil, nil)

	if actual != elevenResults {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 11, elevenResults, actual)
	}

	actual = translator.Translate(key, 4, nil, nil)

	if actual != fourResults {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 4, fourResults, actual)
	}

	actual = translator.Translate(key, 101, nil, nil)

	if actual != results {
		t.Errorf("Should translate '%v' with num '%v' to '%v' but the result is '%v'", key, 101, results, actual)
	}
}
