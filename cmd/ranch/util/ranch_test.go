package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RanchTestSuite struct {
	suite.Suite
}

func (suite *RanchTestSuite) TestRanchValidateConfigOnValidConfig() {
	assert := assert.New(suite.T())
	config := &RanchConfig{
		AppName:   "hello-world",
		ImageName: "hello-world",
	}
	errs := RanchValidateConfig(config)
	assert.Empty(errs)
}

func (suite *RanchTestSuite) TestLinterUrl() {
	assert := assert.New(suite.T())
	Version = "1.0.0"
	assert.Equal(LinterUrl("foobar"), "https://github.com/goodeggs/platform/blob/v1.0.0/LINT_RULES.md#foobar")
}

func (suite *RanchTestSuite) TestRanchValidateConfigOnTooShortCronInterval() {
	assert := assert.New(suite.T())
	crons := map[string]string{
		"foobar": "* * * * ? echo too short",
	}
	config := &RanchConfig{
		AppName:   "hello-world",
		ImageName: "hello-world",
		Cron:      crons,
	}
	errs := RanchValidateConfig(config)
	assert.NotEmpty(errs)
}

func TestRanchTestSuite(t *testing.T) {
	suite.Run(t, new(RanchTestSuite))
}
