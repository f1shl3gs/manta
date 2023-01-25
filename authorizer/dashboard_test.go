package authorizer_test

import (
	"bytes"
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"
	"github.com/f1shl3gs/manta/mock"
	"github.com/f1shl3gs/manta/tests"
)

func TestDashboardService_FindDashboardByID(t *testing.T) {
	testCases := []struct {
		name    string
		service manta.DashboardService

		permission manta.Permission
		id         manta.ID

		want error
	}{
		{
			name: "authorized to access id",
			service: &mock.DashboardService{
				FindDashboardByIDFn: func(ctx context.Context, id manta.ID) (*manta.Dashboard, error) {
					return &manta.Dashboard{
						ID:    id,
						OrgID: 10,
					}, nil
				},
			},
			permission: manta.Permission{
				Action: "read",
				Resource: manta.Resource{
					Type: manta.DashboardsResourceType,
					ID:   tests.IDPtr(1),
				},
			},
			id:   1,
			want: nil,
		},
		{
			name: "unauthorized to access id",
			service: &mock.DashboardService{
				FindDashboardByIDFn: func(ctx context.Context, id manta.ID) (*manta.Dashboard, error) {
					return &manta.Dashboard{
						ID:    id,
						OrgID: 10,
					}, nil
				},
			},
			permission: manta.Permission{
				Action: "read",
				Resource: manta.Resource{
					Type: manta.DashboardsResourceType,
					ID:   tests.IDPtr(2),
				},
			},
			id: 1,
			want: &manta.Error{
				Code: manta.EUnauthorized,
				Msg:  "read:orgs/000000000000000a/dashboards/0000000000000001 is unauthorized",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := authorizer.NewDashboardService(tt.service)

			ctx := context.Background()
			ctx = authorizer.SetAuthorizer(ctx, mock.NewAuthorizer(false, []manta.Permission{tt.permission}))

			_, err := s.FindDashboardByID(ctx, tt.id)
			tests.ErrorsEqual(t, err, tt.want)
		})
	}
}

var dashboardCmpOptions = cmp.Options{
	cmp.Comparer(func(x, y []byte) bool {
		return bytes.Equal(x, y)
	}),
	cmp.Transformer("Sort", func(in []*manta.Dashboard) []*manta.Dashboard {
		out := append([]*manta.Dashboard(nil), in...) // Copy input to avoid mutating it
		sort.Slice(out, func(i, j int) bool {
			return out[i].ID.String() > out[j].ID.String()
		})
		return out
	}),
}

func TestDashboardService_FindDashboards(t *testing.T) {
	type wants struct {
		dashboards []*manta.Dashboard
		err        error
	}

	testCases := []struct {
		name string

		service    manta.DashboardService
		permission manta.Permission

		wants wants
	}{
		{
			name: "authorized to see all dashboards",
			service: &mock.DashboardService{
				FindDashboardsFn: func(ctx context.Context, filter manta.DashboardFilter) ([]*manta.Dashboard, error) {
					return []*manta.Dashboard{
						{
							ID:    1,
							OrgID: 10,
						},
						{
							ID:    2,
							OrgID: 10,
						},
						{
							ID:    3,
							OrgID: 11,
						},
					}, nil
				},
			},
			permission: manta.Permission{
				Action: "read",
				Resource: manta.Resource{
					Type: manta.DashboardsResourceType,
				},
			},
			wants: wants{
				dashboards: []*manta.Dashboard{
					{
						ID:    1,
						OrgID: 10,
					},
					{
						ID:    2,
						OrgID: 10,
					},
					{
						ID:    3,
						OrgID: 11,
					},
				},
			},
		},

		{
			name: "authorized to access a single orgs dashboards",
			service: &mock.DashboardService{
				FindDashboardsFn: func(ctx context.Context, filter manta.DashboardFilter) ([]*manta.Dashboard, error) {
					return []*manta.Dashboard{
						{
							ID:    1,
							OrgID: 10,
						},
						{
							ID:    2,
							OrgID: 10,
						},
						{
							ID:    3,
							OrgID: 11,
						},
					}, nil
				},
			},
			permission: manta.Permission{
				Action: "read",
				Resource: manta.Resource{
					Type:  manta.DashboardsResourceType,
					OrgID: tests.IDPtr(10),
				},
			},
			wants: wants{
				dashboards: []*manta.Dashboard{
					{
						ID:    1,
						OrgID: 10,
					},
					{
						ID:    2,
						OrgID: 10,
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := authorizer.NewDashboardService(tt.service)

			ctx := context.Background()
			ctx = authorizer.SetAuthorizer(ctx, mock.NewAuthorizer(false, []manta.Permission{tt.permission}))

			dashboards, err := s.FindDashboards(ctx, manta.DashboardFilter{})
			tests.ErrorsEqual(t, err, tt.wants.err)

			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}
