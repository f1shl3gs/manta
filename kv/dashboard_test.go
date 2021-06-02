package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/stretchr/testify/require"
)

func initDashboardService(t *testing.T) (manta.DashboardService, func()) {
	svc, closer := NewTestService(t)

	return svc, closer
}

func TestDashboard(t *testing.T) {
	t.Run("create dashboard", func(t *testing.T) {
		svc, closer := NewTestService(t)
		defer closer()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		org := &manta.Organization{
			Name: "test",
			Desc: "test desc",
		}

		err := svc.CreateOrganization(ctx, org)
		require.NoError(t, err)

		d := &manta.Dashboard{
			Name:  "dash",
			Desc:  "dashboard",
			OrgID: org.ID,
			Cells: nil,
		}

		err = svc.CreateDashboard(ctx, d)
		require.NoError(t, err)

		created, err := svc.FindDashboardByID(ctx, d.ID)
		require.NoError(t, err)

		require.Equal(t, d.Name, created.Name)
		require.Equal(t, d.OrgID, created.OrgID)
	})
}
