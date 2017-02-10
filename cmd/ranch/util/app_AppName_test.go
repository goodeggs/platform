package util

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
)

type AppNameTestSuite struct {
	suite.Suite
}

func (suite *AppNameTestSuite) SetupSuite() {
	removeTestFiles()
}

func (suite *AppNameTestSuite) TestRun() {
	assert := assert.New(suite.T())

	// create a junk cobra command for testing
	var testCmd = cobra.Command{
		Use: "junk",
	}
	testCmd.Flags().StringP("filename", "f", "", "config filename (defaults to .ranch.yaml)")

	// get current directory name
	// currentDirectory := "util"
	wd, err := os.Getwd()
	assert.Nil(err)
	currentDirectory := path.Base(wd)

	// test case: use specified app name, regardless of what files exist
	testAppName := "superuberfantastic"
	name, err := _appName(testAppName, &testCmd)
	assert.Equal(testAppName, name, "app name should match specified value")
	assert.Nil(err)

	// test case: .ranch.yaml file does not exist, so use the directory name
	name, err = _appName("", &testCmd)
	assert.Equal(currentDirectory, name, "app name should use directory name")
	assert.Nil(err)

	// create .ranch.yaml file if it 
	os.Create(TEST_RANCHY)

	// test case: no app name specified, should use what is inside .ranch.yaml
	name, err = _appName("", &testCmd)
	assert.Equal("", name, "app name should use directory name when .ranch.yaml is useless")
	assert.Nil(err)

	// Note: this is basically doing a brief test of LoadAppConfig
	value := "name: fantastic\n"
	err = ioutil.WriteFile(TEST_RANCHY, []byte(value), 0644)
	assert.Nil(err)

	// test case: .ranch.yaml file specifies an app name, use that
	name, err = _appName("", &testCmd)
	assert.Equal("fantastic", name, "app name should use directory name when .ranch.yaml is useless")
	assert.Nil(err)

	// test case: no config file passed, scans cwd for .ranch.*.yaml, finds one,
	// so doesn't require you to specify -f
	os.Create(TEST_FILE_1)
	name, err = _appName("", &testCmd)
	assert.Equal("fantastic", name, "file returned should be '.ranch.yaml'")
	assert.Nil(err)

	// test case: no config file passed, scans cwd for .ranch.*.yaml, finds two,
	// so then should error
	os.Create(TEST_FILE_2)
	name, err = _appName("", &testCmd)
	assert.Equal("", name)
	assert.Error(err, "should complain about multiple ranch configs")

	// test case: app name is provided, use even when configs exist
	name, err = _appName(testAppName, &testCmd)
	assert.Equal(testAppName, name, "returned file should be test file 1")
	assert.Nil(err)
}

func (suite *AppNameTestSuite) TestTearDown() {
	removeTestFiles()
}

func TestAppNameTestSuite(t *testing.T) {
	suite.Run(t, new(AppNameTestSuite))
}
