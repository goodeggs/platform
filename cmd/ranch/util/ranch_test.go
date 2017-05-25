package util

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RanchTestSuite struct {
	suite.Suite
}

func createValidRanchConfig() *RanchConfig {
	config := CreateRanchConfig()
	config.AppName = "cool-app"
	config.ImageName = "cool-app"
	return config
}

func (suite *RanchTestSuite) TestRanchLoadRanchConfigDefaultCronMemory() {
	assert := assert.New(suite.T())

	file, err := ioutil.TempFile(os.TempDir(), "ranch_test")
	defer os.Remove(file.Name())

	content := `
name: myapp
`
	if _, err := file.Write([]byte(content)); err != nil {
		assert.Nil(err)
	}

	config, err := LoadRanchConfig(file.Name())
	assert.Nil(err)
	assert.Equal(config.CronMemory, 2048)
}

func (suite *RanchTestSuite) TestRanchLoadRanchConfigExplicitCronMemory() {
	assert := assert.New(suite.T())

	file, err := ioutil.TempFile(os.TempDir(), "ranch_test")
	defer os.Remove(file.Name())

	content := `
name: myapp
cron_memory: 4096
`
	if _, err := file.Write([]byte(content)); err != nil {
		assert.Nil(err)
	}

	config, err := LoadRanchConfig(file.Name())
	assert.Nil(err)
	assert.Equal(config.CronMemory, 4096)
}

func (suite *RanchTestSuite) TestValidateRanchConfigCronMemory() {
	assert := assert.New(suite.T())

	config := createValidRanchConfig()
	config.CronMemory = 3
	errs := RanchValidateConfig(config)
	assert.NotEmpty(errs)

	config.CronMemory = 6144
	errs = RanchValidateConfig(config)
	assert.NotEmpty(errs)

	config.CronMemory = 4096
	errs = RanchValidateConfig(config)
	assert.Empty(errs)
}

func (suite *RanchTestSuite) TestRanchValidateConfigNoCronJobs() {
	assert := assert.New(suite.T())
	config := createValidRanchConfig()
	config.Cron = make(map[string]string)
	errs := RanchValidateConfig(config)
	assert.Empty(errs)
}

func (suite *RanchTestSuite) TestRanchValidateConfigValidCronJobs() {
	assert := assert.New(suite.T())
	crons := map[string]string{
		"one":            "1 * * * ? echo foo",
		"fifteen":        "15 * * * ? echo bar",
		"every10":        "*/10 * * * ? echo bar", // 00,10,20,30,40,50
		"every10onthe3s": "3/10 * * * ? echo bar", // 03,13,23,33,43,53
	}
	config := createValidRanchConfig()
	config.Cron = crons
	errs := RanchValidateConfig(config)
	assert.Empty(errs)
}

func (suite *RanchTestSuite) TestRanchValidateConfigOnTooShortCronInterval() {
	assert := assert.New(suite.T())
	crons := map[string]string{
		"foobar": "* * * * ? echo too short",
	}
	config := createValidRanchConfig()
	config.Cron = crons
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
