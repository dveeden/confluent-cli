// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: repository.go

package mock

import (
	sync "sync"

	github_com_hashicorp_go_version "github.com/hashicorp/go-version"
)

// Repository is a mock of Repository interface
type Repository struct {
	lockGetLatestMajorAndMinorVersion sync.Mutex
	GetLatestMajorAndMinorVersionFunc func(name string, current *github_com_hashicorp_go_version.Version) (*github_com_hashicorp_go_version.Version, *github_com_hashicorp_go_version.Version, error)

	lockGetLatestReleaseNotesVersions sync.Mutex
	GetLatestReleaseNotesVersionsFunc func(name, currentVersion string) (github_com_hashicorp_go_version.Collection, error)

	lockGetAvailableBinaryVersions sync.Mutex
	GetAvailableBinaryVersionsFunc func(name string) (github_com_hashicorp_go_version.Collection, error)

	lockGetAvailableReleaseNotesVersions sync.Mutex
	GetAvailableReleaseNotesVersionsFunc func(name string) (github_com_hashicorp_go_version.Collection, error)

	lockDownloadVersion sync.Mutex
	DownloadVersionFunc func(name, version, downloadDir string) (string, int64, error)

	lockDownloadReleaseNotes sync.Mutex
	DownloadReleaseNotesFunc func(name, version string) (string, error)

	calls struct {
		GetLatestMajorAndMinorVersion []struct {
			Name    string
			Current *github_com_hashicorp_go_version.Version
		}
		GetLatestReleaseNotesVersions []struct {
			Name           string
			CurrentVersion string
		}
		GetAvailableBinaryVersions []struct {
			Name string
		}
		GetAvailableReleaseNotesVersions []struct {
			Name string
		}
		DownloadVersion []struct {
			Name        string
			Version     string
			DownloadDir string
		}
		DownloadReleaseNotes []struct {
			Name    string
			Version string
		}
	}
}

// GetLatestMajorAndMinorVersion mocks base method by wrapping the associated func.
func (m *Repository) GetLatestMajorAndMinorVersion(name string, current *github_com_hashicorp_go_version.Version) (*github_com_hashicorp_go_version.Version, *github_com_hashicorp_go_version.Version, error) {
	m.lockGetLatestMajorAndMinorVersion.Lock()
	defer m.lockGetLatestMajorAndMinorVersion.Unlock()

	if m.GetLatestMajorAndMinorVersionFunc == nil {
		panic("mocker: Repository.GetLatestMajorAndMinorVersionFunc is nil but Repository.GetLatestMajorAndMinorVersion was called.")
	}

	call := struct {
		Name    string
		Current *github_com_hashicorp_go_version.Version
	}{
		Name:    name,
		Current: current,
	}

	m.calls.GetLatestMajorAndMinorVersion = append(m.calls.GetLatestMajorAndMinorVersion, call)

	return m.GetLatestMajorAndMinorVersionFunc(name, current)
}

// GetLatestMajorAndMinorVersionCalled returns true if GetLatestMajorAndMinorVersion was called at least once.
func (m *Repository) GetLatestMajorAndMinorVersionCalled() bool {
	m.lockGetLatestMajorAndMinorVersion.Lock()
	defer m.lockGetLatestMajorAndMinorVersion.Unlock()

	return len(m.calls.GetLatestMajorAndMinorVersion) > 0
}

// GetLatestMajorAndMinorVersionCalls returns the calls made to GetLatestMajorAndMinorVersion.
func (m *Repository) GetLatestMajorAndMinorVersionCalls() []struct {
	Name    string
	Current *github_com_hashicorp_go_version.Version
} {
	m.lockGetLatestMajorAndMinorVersion.Lock()
	defer m.lockGetLatestMajorAndMinorVersion.Unlock()

	return m.calls.GetLatestMajorAndMinorVersion
}

// GetLatestReleaseNotesVersions mocks base method by wrapping the associated func.
func (m *Repository) GetLatestReleaseNotesVersions(name, currentVersion string) (github_com_hashicorp_go_version.Collection, error) {
	m.lockGetLatestReleaseNotesVersions.Lock()
	defer m.lockGetLatestReleaseNotesVersions.Unlock()

	if m.GetLatestReleaseNotesVersionsFunc == nil {
		panic("mocker: Repository.GetLatestReleaseNotesVersionsFunc is nil but Repository.GetLatestReleaseNotesVersions was called.")
	}

	call := struct {
		Name           string
		CurrentVersion string
	}{
		Name:           name,
		CurrentVersion: currentVersion,
	}

	m.calls.GetLatestReleaseNotesVersions = append(m.calls.GetLatestReleaseNotesVersions, call)

	return m.GetLatestReleaseNotesVersionsFunc(name, currentVersion)
}

// GetLatestReleaseNotesVersionsCalled returns true if GetLatestReleaseNotesVersions was called at least once.
func (m *Repository) GetLatestReleaseNotesVersionsCalled() bool {
	m.lockGetLatestReleaseNotesVersions.Lock()
	defer m.lockGetLatestReleaseNotesVersions.Unlock()

	return len(m.calls.GetLatestReleaseNotesVersions) > 0
}

// GetLatestReleaseNotesVersionsCalls returns the calls made to GetLatestReleaseNotesVersions.
func (m *Repository) GetLatestReleaseNotesVersionsCalls() []struct {
	Name           string
	CurrentVersion string
} {
	m.lockGetLatestReleaseNotesVersions.Lock()
	defer m.lockGetLatestReleaseNotesVersions.Unlock()

	return m.calls.GetLatestReleaseNotesVersions
}

// GetAvailableBinaryVersions mocks base method by wrapping the associated func.
func (m *Repository) GetAvailableBinaryVersions(name string) (github_com_hashicorp_go_version.Collection, error) {
	m.lockGetAvailableBinaryVersions.Lock()
	defer m.lockGetAvailableBinaryVersions.Unlock()

	if m.GetAvailableBinaryVersionsFunc == nil {
		panic("mocker: Repository.GetAvailableBinaryVersionsFunc is nil but Repository.GetAvailableBinaryVersions was called.")
	}

	call := struct {
		Name string
	}{
		Name: name,
	}

	m.calls.GetAvailableBinaryVersions = append(m.calls.GetAvailableBinaryVersions, call)

	return m.GetAvailableBinaryVersionsFunc(name)
}

// GetAvailableBinaryVersionsCalled returns true if GetAvailableBinaryVersions was called at least once.
func (m *Repository) GetAvailableBinaryVersionsCalled() bool {
	m.lockGetAvailableBinaryVersions.Lock()
	defer m.lockGetAvailableBinaryVersions.Unlock()

	return len(m.calls.GetAvailableBinaryVersions) > 0
}

// GetAvailableBinaryVersionsCalls returns the calls made to GetAvailableBinaryVersions.
func (m *Repository) GetAvailableBinaryVersionsCalls() []struct {
	Name string
} {
	m.lockGetAvailableBinaryVersions.Lock()
	defer m.lockGetAvailableBinaryVersions.Unlock()

	return m.calls.GetAvailableBinaryVersions
}

// GetAvailableReleaseNotesVersions mocks base method by wrapping the associated func.
func (m *Repository) GetAvailableReleaseNotesVersions(name string) (github_com_hashicorp_go_version.Collection, error) {
	m.lockGetAvailableReleaseNotesVersions.Lock()
	defer m.lockGetAvailableReleaseNotesVersions.Unlock()

	if m.GetAvailableReleaseNotesVersionsFunc == nil {
		panic("mocker: Repository.GetAvailableReleaseNotesVersionsFunc is nil but Repository.GetAvailableReleaseNotesVersions was called.")
	}

	call := struct {
		Name string
	}{
		Name: name,
	}

	m.calls.GetAvailableReleaseNotesVersions = append(m.calls.GetAvailableReleaseNotesVersions, call)

	return m.GetAvailableReleaseNotesVersionsFunc(name)
}

// GetAvailableReleaseNotesVersionsCalled returns true if GetAvailableReleaseNotesVersions was called at least once.
func (m *Repository) GetAvailableReleaseNotesVersionsCalled() bool {
	m.lockGetAvailableReleaseNotesVersions.Lock()
	defer m.lockGetAvailableReleaseNotesVersions.Unlock()

	return len(m.calls.GetAvailableReleaseNotesVersions) > 0
}

// GetAvailableReleaseNotesVersionsCalls returns the calls made to GetAvailableReleaseNotesVersions.
func (m *Repository) GetAvailableReleaseNotesVersionsCalls() []struct {
	Name string
} {
	m.lockGetAvailableReleaseNotesVersions.Lock()
	defer m.lockGetAvailableReleaseNotesVersions.Unlock()

	return m.calls.GetAvailableReleaseNotesVersions
}

// DownloadVersion mocks base method by wrapping the associated func.
func (m *Repository) DownloadVersion(name, version, downloadDir string) (string, int64, error) {
	m.lockDownloadVersion.Lock()
	defer m.lockDownloadVersion.Unlock()

	if m.DownloadVersionFunc == nil {
		panic("mocker: Repository.DownloadVersionFunc is nil but Repository.DownloadVersion was called.")
	}

	call := struct {
		Name        string
		Version     string
		DownloadDir string
	}{
		Name:        name,
		Version:     version,
		DownloadDir: downloadDir,
	}

	m.calls.DownloadVersion = append(m.calls.DownloadVersion, call)

	return m.DownloadVersionFunc(name, version, downloadDir)
}

// DownloadVersionCalled returns true if DownloadVersion was called at least once.
func (m *Repository) DownloadVersionCalled() bool {
	m.lockDownloadVersion.Lock()
	defer m.lockDownloadVersion.Unlock()

	return len(m.calls.DownloadVersion) > 0
}

// DownloadVersionCalls returns the calls made to DownloadVersion.
func (m *Repository) DownloadVersionCalls() []struct {
	Name        string
	Version     string
	DownloadDir string
} {
	m.lockDownloadVersion.Lock()
	defer m.lockDownloadVersion.Unlock()

	return m.calls.DownloadVersion
}

// DownloadReleaseNotes mocks base method by wrapping the associated func.
func (m *Repository) DownloadReleaseNotes(name, version string) (string, error) {
	m.lockDownloadReleaseNotes.Lock()
	defer m.lockDownloadReleaseNotes.Unlock()

	if m.DownloadReleaseNotesFunc == nil {
		panic("mocker: Repository.DownloadReleaseNotesFunc is nil but Repository.DownloadReleaseNotes was called.")
	}

	call := struct {
		Name    string
		Version string
	}{
		Name:    name,
		Version: version,
	}

	m.calls.DownloadReleaseNotes = append(m.calls.DownloadReleaseNotes, call)

	return m.DownloadReleaseNotesFunc(name, version)
}

// DownloadReleaseNotesCalled returns true if DownloadReleaseNotes was called at least once.
func (m *Repository) DownloadReleaseNotesCalled() bool {
	m.lockDownloadReleaseNotes.Lock()
	defer m.lockDownloadReleaseNotes.Unlock()

	return len(m.calls.DownloadReleaseNotes) > 0
}

// DownloadReleaseNotesCalls returns the calls made to DownloadReleaseNotes.
func (m *Repository) DownloadReleaseNotesCalls() []struct {
	Name    string
	Version string
} {
	m.lockDownloadReleaseNotes.Lock()
	defer m.lockDownloadReleaseNotes.Unlock()

	return m.calls.DownloadReleaseNotes
}

// Reset resets the calls made to the mocked methods.
func (m *Repository) Reset() {
	m.lockGetLatestMajorAndMinorVersion.Lock()
	m.calls.GetLatestMajorAndMinorVersion = nil
	m.lockGetLatestMajorAndMinorVersion.Unlock()
	m.lockGetLatestReleaseNotesVersions.Lock()
	m.calls.GetLatestReleaseNotesVersions = nil
	m.lockGetLatestReleaseNotesVersions.Unlock()
	m.lockGetAvailableBinaryVersions.Lock()
	m.calls.GetAvailableBinaryVersions = nil
	m.lockGetAvailableBinaryVersions.Unlock()
	m.lockGetAvailableReleaseNotesVersions.Lock()
	m.calls.GetAvailableReleaseNotesVersions = nil
	m.lockGetAvailableReleaseNotesVersions.Unlock()
	m.lockDownloadVersion.Lock()
	m.calls.DownloadVersion = nil
	m.lockDownloadVersion.Unlock()
	m.lockDownloadReleaseNotes.Lock()
	m.calls.DownloadReleaseNotes = nil
	m.lockDownloadReleaseNotes.Unlock()
}
