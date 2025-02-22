package test

import (
	"strings"

	"github.com/confluentinc/bincover"
)

func (s *CLITestSuite) TestEnvironment() {
	tests := []CLITest{
		// Only log in at the beginning so active env is not reset
		// tt.workflow=true so login is not reset
		{args: "environment list", fixture: "environment/1.golden", login: "cloud"},
		{args: "environment use not-595", fixture: "environment/2.golden"},
		{args: "environment update not-595 --name new-other-name", fixture: "environment/3.golden"},
		{args: "environment list", fixture: "environment/4.golden"},
		{args: "environment list -o json", fixture: "environment/5.golden"},
		{args: "environment list -o yaml", fixture: "environment/6.golden"},
		{args: "environment use non-existent-id", fixture: "environment/7.golden", wantErrCode: 1},
		{args: "environment create saucayyy", fixture: "environment/8.golden"},
		{args: "environment create saucayyy -o json", fixture: "environment/9.golden"},
		{args: "environment create saucayyy -o yaml", fixture: "environment/10.golden"},
		{args: "environment delete not-595 --force", fixture: "environment/11.golden"},
		{args: "environment delete not-595", preCmdFuncs: []bincover.PreCmdFunc{stdinPipeFunc(strings.NewReader("default\n"))}, fixture: "environment/11-prompt.golden"},
		{args: "environment delete env-dne --force", fixture: "environment/12.golden", wantErrCode: 1},
	}

	resetConfiguration(s.T(), false)

	for _, tt := range tests {
		tt.workflow = true
		s.runIntegrationTest(tt)
	}
}

func (s *CLITestSuite) TestEnvironmentDescribe() {
	tests := []CLITest{
		{args: "environment describe env-12345", fixture: "environment/describe.golden"},
		{args: "environment describe env-12345 -o json", fixture: "environment/describe-json.golden"},
	}

	for _, tt := range tests {
		tt.login = "cloud"
		s.runIntegrationTest(tt)
	}
}
