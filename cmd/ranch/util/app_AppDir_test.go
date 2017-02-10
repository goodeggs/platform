package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/cobra"
)

type AppDirTestSuite struct {
	suite.Suite
}

func (suite *AppDirTestSuite) TestRun() {
	assert := assert.New(suite.T())

	// create a junk cobra command for testing
	var testCmd = cobra.Command{
		Use: "junk",
	}

	// TODO: this needs some better coverage, might need to stub out the file system
	_, err := AppDir(&testCmd)
	assert.Nil(err)
}

func TestAppDirTestSuite(t *testing.T) {
	suite.Run(t, new(AppDirTestSuite))
}
