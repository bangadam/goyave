package lang

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"goyave.dev/goyave/v4/config"
)

type LangTestSuite struct {
	suite.Suite
}

func loadTestLang(lang string) {
	_, filename, _, _ := runtime.Caller(1)
	load(lang, path.Dir(filename)+"/../resources/lang/en-US")
}

func (suite *LangTestSuite) SetupSuite() {
	LoadDefault()
	LoadAllAvailableLanguages()

	if err := config.LoadFrom("../config.test.json"); err != nil {
		suite.FailNow(err.Error())
	}
	config.Set("app.defaultLanguage", "en-US")
}

func (suite *LangTestSuite) TestLang() {
	suite.Equal("email address", Get("en-US", "validation.fields.email"))
	suite.Equal("The :field is required.", Get("en-US", "validation.rules.required"))
	suite.Equal("Malformed request", Get("en-US", "malformed-request"))
	suite.Equal("Invalid credentials.", Get("en-US", "auth.invalid-credentials"))
	suite.Equal("doesn't.exist", Get("en-US", "doesn't.exist"))
	suite.Equal("doesn'texist", Get("en-US", "doesn'texist"))
	suite.Equal("validation.doesn't.exist", Get("en-US", "validation.doesn't.exist"))
	suite.Equal("validation.rules", Get("en-US", "validation.rules"))
	suite.Equal("validation.rules.doesn't.exist", Get("en-US", "validation.rules.doesn't.exist"))
	suite.Equal("validation.fields.doesn't", Get("en-US", "validation.fields.doesn't"))
	suite.Equal("validation.fields.doesn.t.", Get("en-US", "validation.fields.doesn.t."))

	suite.Equal("validation.fields", Get("en-US", "validation.fields"))
	suite.Equal("doesn't.exist", Get("not a language", "doesn't.exist"))
}

func (suite *LangTestSuite) TestDetectLanguage() {
	loadTestLang("fr-FR")
	loadTestLang("fr-FR") // Merge existing

	suite.Equal("en-US", DetectLanguage("en"))
	suite.Equal("en-US", DetectLanguage("en-US, fr"))
	suite.Equal("fr-FR", DetectLanguage("fr-FR"))
	suite.Equal("fr-FR", DetectLanguage("fr"))
	suite.Equal("en-US", DetectLanguage("fr, en-US"))
	suite.Equal("fr-FR", DetectLanguage("fr-FR, en-US"))
	suite.Equal("fr-FR", DetectLanguage("fr, en-US;q=0.9"))
	suite.Equal("en-US", DetectLanguage("en, fr-FR;q=0.9"))
	suite.Equal("en-US", DetectLanguage("*"))
	suite.Equal("en-US", DetectLanguage("notalang"))

	langs := GetAvailableLanguages()
	suite.Equal(2, len(langs))
	suite.Contains(langs, "en-US")
	suite.Contains(langs, "fr-FR")
}

func (suite *LangTestSuite) TestLoad() {
	suite.Panics(func() {
		Load("notalanguagedir", "notalanguagepath")
	})

	Load("en-US", "../resources/lang/en-US") // Is an override
	suite.Equal("rule override", languages["en-US"].validation.rules["required"])

	dest := map[string]string{}
	err := readLangFile("../resources/lang/invalid.json", &dest)
	suite.NotNil(err)

	// Ensure default lang is not changed
	suite.Equal("The :field is required.", enUS.validation.rules["required"])

}

func (suite *LangTestSuite) TestMerge() {
	dst := &Language{
		lines: map[string]string{"line": "line 1"},
		validation: validationLines{
			rules: map[string]string{},
			fields: map[string]string{
				"test": "test field",
			},
		},
	}
	src := &Language{
		lines: map[string]string{"other": "line 2"},
		validation: validationLines{
			rules: map[string]string{},
			fields: map[string]string{
				"email": "email address",
				"test":  "test field override",
			},
		},
	}

	mergeLang(dst, src)
	suite.Equal("line 1", dst.lines["line"])
	suite.Equal("line 2", dst.lines["other"])

	suite.Equal("email address", dst.validation.fields["email"])

	suite.Equal("test field override", dst.validation.fields["test"])
}

func (suite *LangTestSuite) TestPlaceholders() {
	suite.Equal("Greetings, Kevin", convertEmptyLine("greetings", "Greetings, :username", []string{":username", "Kevin"}))
	suite.Equal("Greetings, Kevin, today is Monday", convertEmptyLine("greetings", "Greetings, :username, today is :today", []string{":username", "Kevin", ":today", "Monday"}))
	suite.Equal("Greetings, Kevin, today is :today", convertEmptyLine("greetings", "Greetings, :username, today is :today", []string{":username", "Kevin", ":today"}))
}

func (suite *LangTestSuite) TestSetDefault() {
	SetDefaultLine("test-line", "It's sunny today")
	suite.Equal("It's sunny today", enUS.lines["test-line"])
	delete(enUS.lines, "test-line")

	SetDefaultValidationRule("test-validation-rules", "It's sunny today")
	suite.Equal("It's sunny today", enUS.validation.rules["test-validation-rules"])
	delete(enUS.validation.rules, "test-validation-rules")

	SetDefaultFieldName("test-field-name", "Sun")
	suite.Equal("Sun", enUS.validation.fields["test-field-name"])
	delete(enUS.validation.fields, "test-field-name")

	// Test no override
	enUS.validation.fields["test-field"] = "test"
	SetDefaultFieldName("test-field", "Sun")
	suite.Equal("Sun", enUS.validation.fields["test-field"])
	delete(enUS.validation.fields, "test-field")
}

func (suite *LangTestSuite) TearDownAllSuite() {
	languages = map[string]*Language{}
}

func TestLangTestSuite(t *testing.T) {
	suite.Run(t, new(LangTestSuite))
}
