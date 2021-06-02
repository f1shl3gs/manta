package kv_test

import (
	"context"
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/require"
)

func initCheckService(t *testing.T) (manta.ID, manta.CheckService, manta.TaskService, func()) {
	svc, closer := NewTestService(t)
	orgID := CreateDefaultOrg(t, svc)
	return orgID, svc, svc, closer
}

func TestCheck(t *testing.T) {
	t.Run("create check", func(t *testing.T) {
		orgID, checkService, taskService, closer := initCheckService(t)
		defer closer()

		check := &manta.Check{
			Name:   "foo",
			Desc:   "foo desc",
			OrgID:  orgID,
			Expr:   "rate(node_cpu_seconds_total[1m])",
			Status: "active",
			Cron:   "@every 1m",
			Conditions: []manta.Condition{
				{
					Status:  "warning",
					Pending: 0,
					Threshold: manta.Threshold{
						Type:  "lt",
						Value: 90.0,
					},
				},
			},
		}

		err := checkService.CreateCheck(context.Background(), check)
		require.NoError(t, err)

		tasks, err := taskService.FindTasks(context.Background(), manta.TaskFilter{OwnerID: &check.ID})
		require.NoError(t, err)

		require.Equal(t, 1, len(tasks))
		task := tasks[0]
		require.Equal(t, check.ID, task.OwnerID)
		require.Equal(t, check.OrgID, task.OrgID)
	})

	t.Run("find checks", func(t *testing.T) {

	})
}
