package task

import (
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

const (
	TaskOptOrgID                = "OrgID"
	TaskOptStartedAt            = "StartedAt"
	TaskOptAlwaysRetryOnFailure = "AlwaysRetryOnFailure"
)

type TaskOpt struct {
	name string
	args []any
	fn   func(*apigen.Task) error
}

func OrgID(orgID int32) TaskOpt {
	return TaskOpt{
		name: "OrgID",
		args: []any{orgID},
		fn: func(t *apigen.Task) error {
			t.Attributes.OrgID = &orgID
			return nil
		},
	}
}

func StartedAt(startedAt time.Time) TaskOpt {
	return TaskOpt{
		name: "StartedAt",
		args: []any{startedAt},
		fn: func(t *apigen.Task) error {
			t.StartedAt = &startedAt
			return nil
		},
	}
}

func AlwaysRetryOnFailure(interval string) TaskOpt {
	return TaskOpt{
		name: "AlwaysRetryOnFailure",
		args: []any{interval},
		fn: func(t *apigen.Task) error {
			_, err := time.ParseDuration(interval)
			if err != nil {
				return errors.Wrap(err, "failed to parse retry interval")
			}
			if t.Attributes.RetryPolicy == nil {
				t.Attributes.RetryPolicy = &apigen.TaskRetryPolicy{}
			}
			t.Attributes.RetryPolicy.AlwaysRetryOnFailure = utils.Ptr(true)
			t.Attributes.RetryPolicy.Interval = interval
			return nil
		},
	}
}
