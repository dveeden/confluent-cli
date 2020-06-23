// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: analytics.go

package mock

import (
	sync "sync"

	github_com_confluentinc_cli_internal_pkg_analytics "github.com/confluentinc/cli/internal/pkg/analytics"
	github_com_spf13_cobra "github.com/spf13/cobra"
)

// AnalyticsClient is a mock of Client interface
type AnalyticsClient struct {
	lockSetStartTime sync.Mutex
	SetStartTimeFunc func()

	lockTrackCommand sync.Mutex
	TrackCommandFunc func(cmd *github_com_spf13_cobra.Command, args []string)

	lockSetCommandType sync.Mutex
	SetCommandTypeFunc func(commandType github_com_confluentinc_cli_internal_pkg_analytics.CommandType)

	lockSessionTimedOut sync.Mutex
	SessionTimedOutFunc func() error

	lockSendCommandAnalytics sync.Mutex
	SendCommandAnalyticsFunc func(cmd *github_com_spf13_cobra.Command, args []string, cmdExecutionError error) error

	lockClose sync.Mutex
	CloseFunc func() error

	lockSetSpecialProperty sync.Mutex
	SetSpecialPropertyFunc func(propertiesKey string, value interface{})

	calls struct {
		SetStartTime []struct {
		}
		TrackCommand []struct {
			Cmd  *github_com_spf13_cobra.Command
			Args []string
		}
		SetCommandType []struct {
			CommandType github_com_confluentinc_cli_internal_pkg_analytics.CommandType
		}
		SessionTimedOut []struct {
		}
		SendCommandAnalytics []struct {
			Cmd               *github_com_spf13_cobra.Command
			Args              []string
			CmdExecutionError error
		}
		Close []struct {
		}
		SetSpecialProperty []struct {
			PropertiesKey string
			Value         interface{}
		}
	}
}

// SetStartTime mocks base method by wrapping the associated func.
func (m *AnalyticsClient) SetStartTime() {
	m.lockSetStartTime.Lock()
	defer m.lockSetStartTime.Unlock()

	if m.SetStartTimeFunc == nil {
		panic("mocker: AnalyticsClient.SetStartTimeFunc is nil but AnalyticsClient.SetStartTime was called.")
	}

	call := struct {
	}{}

	m.calls.SetStartTime = append(m.calls.SetStartTime, call)

	m.SetStartTimeFunc()
}

// SetStartTimeCalled returns true if SetStartTime was called at least once.
func (m *AnalyticsClient) SetStartTimeCalled() bool {
	m.lockSetStartTime.Lock()
	defer m.lockSetStartTime.Unlock()

	return len(m.calls.SetStartTime) > 0
}

// SetStartTimeCalls returns the calls made to SetStartTime.
func (m *AnalyticsClient) SetStartTimeCalls() []struct {
} {
	m.lockSetStartTime.Lock()
	defer m.lockSetStartTime.Unlock()

	return m.calls.SetStartTime
}

// TrackCommand mocks base method by wrapping the associated func.
func (m *AnalyticsClient) TrackCommand(cmd *github_com_spf13_cobra.Command, args []string) {
	m.lockTrackCommand.Lock()
	defer m.lockTrackCommand.Unlock()

	if m.TrackCommandFunc == nil {
		panic("mocker: AnalyticsClient.TrackCommandFunc is nil but AnalyticsClient.TrackCommand was called.")
	}

	call := struct {
		Cmd  *github_com_spf13_cobra.Command
		Args []string
	}{
		Cmd:  cmd,
		Args: args,
	}

	m.calls.TrackCommand = append(m.calls.TrackCommand, call)

	m.TrackCommandFunc(cmd, args)
}

// TrackCommandCalled returns true if TrackCommand was called at least once.
func (m *AnalyticsClient) TrackCommandCalled() bool {
	m.lockTrackCommand.Lock()
	defer m.lockTrackCommand.Unlock()

	return len(m.calls.TrackCommand) > 0
}

// TrackCommandCalls returns the calls made to TrackCommand.
func (m *AnalyticsClient) TrackCommandCalls() []struct {
	Cmd  *github_com_spf13_cobra.Command
	Args []string
} {
	m.lockTrackCommand.Lock()
	defer m.lockTrackCommand.Unlock()

	return m.calls.TrackCommand
}

// SetCommandType mocks base method by wrapping the associated func.
func (m *AnalyticsClient) SetCommandType(commandType github_com_confluentinc_cli_internal_pkg_analytics.CommandType) {
	m.lockSetCommandType.Lock()
	defer m.lockSetCommandType.Unlock()

	if m.SetCommandTypeFunc == nil {
		panic("mocker: AnalyticsClient.SetCommandTypeFunc is nil but AnalyticsClient.SetCommandType was called.")
	}

	call := struct {
		CommandType github_com_confluentinc_cli_internal_pkg_analytics.CommandType
	}{
		CommandType: commandType,
	}

	m.calls.SetCommandType = append(m.calls.SetCommandType, call)

	m.SetCommandTypeFunc(commandType)
}

// SetCommandTypeCalled returns true if SetCommandType was called at least once.
func (m *AnalyticsClient) SetCommandTypeCalled() bool {
	m.lockSetCommandType.Lock()
	defer m.lockSetCommandType.Unlock()

	return len(m.calls.SetCommandType) > 0
}

// SetCommandTypeCalls returns the calls made to SetCommandType.
func (m *AnalyticsClient) SetCommandTypeCalls() []struct {
	CommandType github_com_confluentinc_cli_internal_pkg_analytics.CommandType
} {
	m.lockSetCommandType.Lock()
	defer m.lockSetCommandType.Unlock()

	return m.calls.SetCommandType
}

// SessionTimedOut mocks base method by wrapping the associated func.
func (m *AnalyticsClient) SessionTimedOut() error {
	m.lockSessionTimedOut.Lock()
	defer m.lockSessionTimedOut.Unlock()

	if m.SessionTimedOutFunc == nil {
		panic("mocker: AnalyticsClient.SessionTimedOutFunc is nil but AnalyticsClient.SessionTimedOut was called.")
	}

	call := struct {
	}{}

	m.calls.SessionTimedOut = append(m.calls.SessionTimedOut, call)

	return m.SessionTimedOutFunc()
}

// SessionTimedOutCalled returns true if SessionTimedOut was called at least once.
func (m *AnalyticsClient) SessionTimedOutCalled() bool {
	m.lockSessionTimedOut.Lock()
	defer m.lockSessionTimedOut.Unlock()

	return len(m.calls.SessionTimedOut) > 0
}

// SessionTimedOutCalls returns the calls made to SessionTimedOut.
func (m *AnalyticsClient) SessionTimedOutCalls() []struct {
} {
	m.lockSessionTimedOut.Lock()
	defer m.lockSessionTimedOut.Unlock()

	return m.calls.SessionTimedOut
}

// SendCommandAnalytics mocks base method by wrapping the associated func.
func (m *AnalyticsClient) SendCommandAnalytics(cmd *github_com_spf13_cobra.Command, args []string, cmdExecutionError error) error {
	m.lockSendCommandAnalytics.Lock()
	defer m.lockSendCommandAnalytics.Unlock()

	if m.SendCommandAnalyticsFunc == nil {
		panic("mocker: AnalyticsClient.SendCommandAnalyticsFunc is nil but AnalyticsClient.SendCommandAnalytics was called.")
	}

	call := struct {
		Cmd               *github_com_spf13_cobra.Command
		Args              []string
		CmdExecutionError error
	}{
		Cmd:               cmd,
		Args:              args,
		CmdExecutionError: cmdExecutionError,
	}

	m.calls.SendCommandAnalytics = append(m.calls.SendCommandAnalytics, call)

	return m.SendCommandAnalyticsFunc(cmd, args, cmdExecutionError)
}

// SendCommandAnalyticsCalled returns true if SendCommandAnalytics was called at least once.
func (m *AnalyticsClient) SendCommandAnalyticsCalled() bool {
	m.lockSendCommandAnalytics.Lock()
	defer m.lockSendCommandAnalytics.Unlock()

	return len(m.calls.SendCommandAnalytics) > 0
}

// SendCommandAnalyticsCalls returns the calls made to SendCommandAnalytics.
func (m *AnalyticsClient) SendCommandAnalyticsCalls() []struct {
	Cmd               *github_com_spf13_cobra.Command
	Args              []string
	CmdExecutionError error
} {
	m.lockSendCommandAnalytics.Lock()
	defer m.lockSendCommandAnalytics.Unlock()

	return m.calls.SendCommandAnalytics
}

// Close mocks base method by wrapping the associated func.
func (m *AnalyticsClient) Close() error {
	m.lockClose.Lock()
	defer m.lockClose.Unlock()

	if m.CloseFunc == nil {
		panic("mocker: AnalyticsClient.CloseFunc is nil but AnalyticsClient.Close was called.")
	}

	call := struct {
	}{}

	m.calls.Close = append(m.calls.Close, call)

	return m.CloseFunc()
}

// CloseCalled returns true if Close was called at least once.
func (m *AnalyticsClient) CloseCalled() bool {
	m.lockClose.Lock()
	defer m.lockClose.Unlock()

	return len(m.calls.Close) > 0
}

// CloseCalls returns the calls made to Close.
func (m *AnalyticsClient) CloseCalls() []struct {
} {
	m.lockClose.Lock()
	defer m.lockClose.Unlock()

	return m.calls.Close
}

// SetSpecialProperty mocks base method by wrapping the associated func.
func (m *AnalyticsClient) SetSpecialProperty(propertiesKey string, value interface{}) {
	m.lockSetSpecialProperty.Lock()
	defer m.lockSetSpecialProperty.Unlock()

	if m.SetSpecialPropertyFunc == nil {
		panic("mocker: AnalyticsClient.SetSpecialPropertyFunc is nil but AnalyticsClient.SetSpecialProperty was called.")
	}

	call := struct {
		PropertiesKey string
		Value         interface{}
	}{
		PropertiesKey: propertiesKey,
		Value:         value,
	}

	m.calls.SetSpecialProperty = append(m.calls.SetSpecialProperty, call)

	m.SetSpecialPropertyFunc(propertiesKey, value)
}

// SetSpecialPropertyCalled returns true if SetSpecialProperty was called at least once.
func (m *AnalyticsClient) SetSpecialPropertyCalled() bool {
	m.lockSetSpecialProperty.Lock()
	defer m.lockSetSpecialProperty.Unlock()

	return len(m.calls.SetSpecialProperty) > 0
}

// SetSpecialPropertyCalls returns the calls made to SetSpecialProperty.
func (m *AnalyticsClient) SetSpecialPropertyCalls() []struct {
	PropertiesKey string
	Value         interface{}
} {
	m.lockSetSpecialProperty.Lock()
	defer m.lockSetSpecialProperty.Unlock()

	return m.calls.SetSpecialProperty
}

// Reset resets the calls made to the mocked methods.
func (m *AnalyticsClient) Reset() {
	m.lockSetStartTime.Lock()
	m.calls.SetStartTime = nil
	m.lockSetStartTime.Unlock()
	m.lockTrackCommand.Lock()
	m.calls.TrackCommand = nil
	m.lockTrackCommand.Unlock()
	m.lockSetCommandType.Lock()
	m.calls.SetCommandType = nil
	m.lockSetCommandType.Unlock()
	m.lockSessionTimedOut.Lock()
	m.calls.SessionTimedOut = nil
	m.lockSessionTimedOut.Unlock()
	m.lockSendCommandAnalytics.Lock()
	m.calls.SendCommandAnalytics = nil
	m.lockSendCommandAnalytics.Unlock()
	m.lockClose.Lock()
	m.calls.Close = nil
	m.lockClose.Unlock()
	m.lockSetSpecialProperty.Lock()
	m.calls.SetSpecialProperty = nil
	m.lockSetSpecialProperty.Unlock()
}
