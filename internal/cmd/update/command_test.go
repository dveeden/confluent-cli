package update

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/confluentinc/cli/internal/pkg/log"
	updatemock "github.com/confluentinc/cli/internal/pkg/update/mock"
	"github.com/confluentinc/cli/internal/pkg/version"
	"github.com/confluentinc/cli/mock"
)

func TestGetReleaseNotes_MultipleReleaseNotes(t *testing.T) {
	client := &updatemock.Client{
		GetLatestReleaseNotesFunc: func(_, _ string) (string, []string, error) {
			notes := []string{
				"v0.1.0 changes\n",
				"v1.0.0 changes\n",
			}
			return "1.0.0", notes, nil
		},
	}

	c := &command{
		client:          client,
		logger:          log.New(),
		analyticsClient: mock.NewDummyAnalyticsMock(),
		version:         &version.Version{Version: "0.0.0"},
	}

	require.Equal(t, "v0.1.0 changes\n\nv1.0.0 changes\n", c.getReleaseNotes("confluent", "1.0.0"))
}
