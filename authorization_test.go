package manta_test

import (
	"testing"

	"github.com/f1shl3gs/manta"
	testutil "github.com/f1shl3gs/manta/tests"
)

func TestAuthorizer_PermissionAllowed(t *testing.T) {
	tests := []struct {
		name        string
		permission  manta.Permission
		permissions []manta.Permission
		allowed     bool
	}{
		{
			name: "global permission",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type: manta.ChecksResourceType,
					},
				},
			},
			allowed: true,
		},
		{
			name: "bad org id in permission",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(0),
					ID:    testutil.IDPtr(0),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(1),
					},
				},
			},
			allowed: false,
		},
		{
			name: "bad resource id in permission",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(0),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(1),
					},
				},
			},
			allowed: false,
		},
		{
			name: "bad resource id in permissions",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(0),
					},
				},
			},
			allowed: false,
		},
		{
			name: "matching action resource and ID",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(1),
					},
				},
			},
			allowed: true,
		},
		{
			name: "matching action resource with total",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
					},
				},
			},
			allowed: true,
		},
		{
			name: "matching action resource no ID",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
					},
				},
			},
			allowed: true,
		},
		{
			name: "matching action resource differing ID",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(2),
					},
				},
			},
			allowed: false,
		},
		{
			name: "differing action same resource",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.ReadAction,
					Resource: manta.Resource{
						Type:  manta.ChecksResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(1),
					},
				},
			},
			allowed: false,
		},
		{
			name: "same action differing resource",
			permission: manta.Permission{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    testutil.IDPtr(1),
				},
			},
			permissions: []manta.Permission{
				{
					Action: manta.WriteAction,
					Resource: manta.Resource{
						Type:  manta.DashboardsResourceType,
						OrgID: testutil.IDPtr(1),
						ID:    testutil.IDPtr(1),
					},
				},
			},
			allowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := manta.PermissionAllowed(tt.permission, tt.permissions)
			if allowed != tt.allowed {
				t.Errorf("got allowed = %v, expected allowed = %v", allowed, tt.allowed)
			}
		})
	}
}

func TestPermission_Valid(t *testing.T) {
	type fields struct {
		Action   manta.Action
		Resource manta.Resource
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid bucket permission with ID",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					ID:    validID(),
					OrgID: testutil.IDPtr(1),
				},
			},
		},
		{
			name: "valid bucket permission with nil ID",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					ID:    nil,
					OrgID: testutil.IDPtr(1),
				},
			},
		},
		{
			name: "invalid bucket permission with an invalid ID",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					ID:    func() *manta.ID { id := manta.InvalidID(); return &id }(),
					OrgID: testutil.IDPtr(1),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid permission without an action",
			fields: fields{
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid permission without a resource",
			fields: fields{
				Action: manta.WriteAction,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &manta.Permission{
				Action:   tt.fields.Action,
				Resource: tt.fields.Resource,
			}
			if err := p.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Permission.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionAllResources_Valid(t *testing.T) {

	for _, rt := range manta.AllResourceTypes {
		p := &manta.Permission{
			Action: manta.WriteAction,
			Resource: manta.Resource{
				Type: rt,
				ID:   testutil.IDPtr(1),
			},
		}

		if err := p.Valid(); err != nil {
			t.Errorf("PermissionAllResources.Valid() error = %v", err)
		}
	}
}

func TestPermissionAllActions(t *testing.T) {
	var actions = []manta.Action{
		manta.ReadAction,
		manta.WriteAction,
	}

	for _, a := range actions {
		p := &manta.Permission{
			Action: a,
			Resource: manta.Resource{
				Type:  manta.TasksResourceType,
				OrgID: testutil.IDPtr(1),
			},
		}

		if err := p.Valid(); err != nil {
			t.Errorf("PermissionAllActions.Valid() error = %v", err)
		}
	}
}

func TestPermission_String(t *testing.T) {
	type fields struct {
		Action   manta.Action
		Resource manta.Resource
		Name     *string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid permission with no id",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
				},
			},
			want: `write:orgs/0000000000000001/checks`,
		},
		{
			name: "valid permission with an id",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type:  manta.ChecksResourceType,
					OrgID: testutil.IDPtr(1),
					ID:    validID(),
				},
			},
            want: `write:orgs/0000000000000001/checks/0000000000000064`,
		},
		{
			name: "valid permission with no id or org id",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type: manta.ChecksResourceType,
				},
			},
            want: `write:checks`,
		},
		{
			name: "valid permission with no org id",
			fields: fields{
				Action: manta.WriteAction,
				Resource: manta.Resource{
					Type: manta.ChecksResourceType,
					ID:   testutil.IDPtr(1),
				},
			},
            want: `write:checks/0000000000000001`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := manta.Permission{
				Action:   tt.fields.Action,
				Resource: tt.fields.Resource,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("Permission.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func validID() *manta.ID {
	id := manta.ID(100)
	return &id
}
