package searchtest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/mattermost/mattermost-server/v5/utils"
)

const (
	ENGINE_ALL           = "all"
	ENGINE_MYSQL         = "mysql"
	ENGINE_POSTGRES      = "postgres"
	ENGINE_ELASTICSEARCH = "elasticsearch"
)

type SearchTestEngine struct {
	Driver     string
	BeforeTest func(*testing.T, store.Store)
	AfterTest  func(*testing.T, store.Store)
}

type searchTest struct {
	Name string
	Fn   func(*testing.T, *SearchTestHelper)
	Tags []string
}

func filterTestsByTag(tests []searchTest, tags ...string) []searchTest {
	filteredTests := []searchTest{}
	for _, test := range tests {
		if utils.StringInSlice(ENGINE_ALL, test.Tags) {
			filteredTests = append(filteredTests, test)
			continue
		}
		for _, tag := range tags {
			if utils.StringInSlice(tag, test.Tags) {
				filteredTests = append(filteredTests, test)
				break
			}
		}
	}

	return filteredTests
}

func runTestSearch(t *testing.T, testEngine *SearchTestEngine, tests []searchTest, th *SearchTestHelper) {
	filteredTests := filterTestsByTag(tests, testEngine.Driver)

	for _, test := range filteredTests {
		if testing.Short() {
			t.Skip("Skipping advanced search test")
			continue
		}

		if testEngine.BeforeTest != nil {
			testEngine.BeforeTest(t, th.Store)
		}
		t.Run(test.Name, func(t *testing.T) { test.Fn(t, th) })
		if testEngine.AfterTest != nil {
			testEngine.AfterTest(t, th.Store)
		}
	}
}
