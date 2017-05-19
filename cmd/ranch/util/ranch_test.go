package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RanchTestSuite struct {
	suite.Suite
}

func (suite *RanchTestSuite) TestRanchValidateConfigNoCronJobs() {
	assert := assert.New(suite.T())
	config := &RanchConfig{
		AppName:   "hello-world",
		ImageName: "hello-world",
	}
	errs := RanchValidateConfig(config)
	assert.Empty(errs)
}

func (suite *RanchTestSuite) TestRanchValidateConfigValidCronJobs() {
	assert := assert.New(suite.T())
	crons := map[string]string{
		"one": "1 * * * ? echo foo",
		"fifteen": "15 * * * ? echo bar",
		"every10": "*/10 * * * ? echo bar", // 00,10,20,30,40,50
		"every10onthe3s": "3/10 * * * ? echo bar", // 03,13,23,33,43,53
	}
	config := &RanchConfig{
		AppName:   "hello-world",
		ImageName: "hello-world",
		Cron:      crons,
	}
	errs := RanchValidateConfig(config)
	assert.Empty(errs)
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

func (suite *RanchTestSuite) TestLinterUrl() {
	assert := assert.New(suite.T())
	Version = "1.0.0"
	assert.Equal(LinterUrl("foobar"), "https://github.com/goodeggs/platform/blob/v1.0.0/cmd/ranch/LINT_RULES.md#foobar")
}

func TestRanchTestSuite(t *testing.T) {
	suite.Run(t, new(RanchTestSuite))
}
