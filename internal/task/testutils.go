package task

import (
	"fmt"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	gomock "go.uber.org/mock/gomock"
)

type TaskOptMatcher struct {
	opt TaskOpt
}

func (m *TaskOptMatcher) Matches(x interface{}) bool {
	opt, ok := x.(TaskOpt)
	if !ok {
		return false
	}
	task := &apigen.Task{}
	if err := opt.fn(task); err != nil {
		return false
	}
	switch opt.name {
	case TaskOptOrgID:
		return task.Attributes.OrgID != nil && *task.Attributes.OrgID == opt.args[0].(int32)
	case TaskOptStartedAt:
		return task.StartedAt != nil && task.StartedAt.Equal(opt.args[0].(time.Time))
	case TaskOptAlwaysRetryOnFailure:
		return task.Attributes.RetryPolicy != nil &&
			task.Attributes.RetryPolicy.AlwaysRetryOnFailure != nil &&
			*task.Attributes.RetryPolicy.AlwaysRetryOnFailure &&
			task.Attributes.RetryPolicy.Interval == opt.args[0].(string)
	}
	return false
}

func (m *TaskOptMatcher) String() string {
	return fmt.Sprintf("is %v", m.opt.name)
}

func MatchTaskOpt(opt TaskOpt) gomock.Matcher {
	return &TaskOptMatcher{opt}
}
