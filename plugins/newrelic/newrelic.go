//go:build newrelic
// +build newrelic

package newrelic

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/shoplineapp/go-app/plugins"
	"github.com/sirupsen/logrus"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewNewrelicAgent)
}

var configOptions []newrelic.ConfigOption

type NewrelicAgent struct {
	app *newrelic.Application
}

func (a NewrelicAgent) App() *newrelic.Application {
	return a.app
}

func Configure(configs ...newrelic.ConfigOption) {
	configOptions = configs
}

func NewNewrelicAgent() *NewrelicAgent {
	a, err := newrelic.NewApplication(configOptions...)
	if err != nil {
		logrus.Error("Unable to load Newrelic application")
	}
	return &NewrelicAgent{app: a}
}
