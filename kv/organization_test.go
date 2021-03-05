package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/kv/migration"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestOrganization(t *testing.T) {
	store, closer := NewTestBolt(t, true)
	defer closer()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc := kv.NewService(zaptest.NewLogger(t), store)
	err := migration.Initial(ctx, store)
	require.NoError(t, err)

	err = svc.CreateOrganization(ctx, &manta.Organization{
		Name: "foo",
		Desc: "foo",
	})
	require.NoError(t, err)

	err = svc.CreateOrganization(ctx, &manta.Organization{
		Name: "bar",
		Desc: "bar",
	})
	require.NoError(t, err)

	t.Run("find all organizations", func(t *testing.T) {
		orgs, _, err := svc.FindOrganizations(ctx, manta.OrganizationFilter{})
		require.NoError(t, err)
		require.Equal(t, 2, len(orgs))
		require.Equal(t, "foo", orgs[0].Name)
		require.Equal(t, "bar", orgs[1].Name)
	})

	t.Run("find by name", func(t *testing.T) {
		name := "foo"
		org, err := svc.FindOrganization(ctx, manta.OrganizationFilter{
			Name: &name,
		})
		require.NoError(t, err)
		require.Equal(t, name, org.Name)
	})

	t.Run("name conflict", func(t *testing.T) {
		err = svc.CreateOrganization(ctx, &manta.Organization{
			Name: "foo",
			Desc: "foo",
		})
		require.Equal(t, kv.ErrKeyConflict, err)
	})
}
