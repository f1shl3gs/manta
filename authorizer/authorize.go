package authorizer

import (
	"context"
	"fmt"

	"github.com/f1shl3gs/manta"
)

func isAllowedAll(authorizer manta.Authorizer, permissions []manta.Permission) error {
	pset, err := authorizer.PermissionSet()
	if err != nil {
		return err
	}

	for _, p := range permissions {
		if !pset.Allowed(p) {
            return &manta.Error{
                Code: manta.EUnauthorized,
				Msg:  fmt.Sprintf("%s is unauthorized", p),
			}
		}
	}

	return nil
}

func isAllowed(authorizer manta.Authorizer, p manta.Permission) error {
	return isAllowedAll(authorizer, []manta.Permission{p})
}

func authorize(ctx context.Context, action manta.Action, rt manta.ResourceType, rid, oid *manta.ID) (manta.Authorizer, manta.Permission, error) {
	var (
		p   *manta.Permission
		err error
	)

	if rid != nil && oid != nil {
		p, err = manta.NewPermissionAtID(*rid, action, rt, *oid)
	} else if rid != nil {
		p, err = manta.NewResourcePermission(action, rt, *rid)
	} else if oid != nil {
		p, err = manta.NewPermission(action, rt, *oid)
	} else {
		p, err = manta.NewGlobalPermission(action, rt)
	}

	if err != nil {
		return nil, manta.Permission{}, err
	}

	auth, err := FromContext(ctx)
	if err != nil {
		return nil, manta.Permission{}, err
	}

	return auth, *p, isAllowed(auth, *p)
}

// AuthorizeRead authorizes the user in the context to read the specified resource (identified by its type, ID, and orgID).
// NOTE: authorization will pass even if the user only has permissions for the resource type and organization ID only.
func authorizeRead(ctx context.Context, rt manta.ResourceType, rid, oid manta.ID) (manta.Authorizer, manta.Permission, error) {
	return authorize(ctx, manta.ReadAction, rt, &rid, &oid)
}

// AuthorizeWrite authorizes the user in the context to write the specified resource (identified by its type, ID, and orgID).
// NOTE: authorization will pass even if the user only has permissions for the resource type and organization ID only.
func authorizeWrite(ctx context.Context, rt manta.ResourceType, rid, oid manta.ID) (manta.Authorizer, manta.Permission, error) {
	return authorize(ctx, manta.WriteAction, rt, &rid, &oid)
}

// authorizeCreate authorizes a user to create a resource of the given type for the given org
func authorizeCreate(ctx context.Context, rt manta.ResourceType, oid manta.ID) (manta.Authorizer, manta.Permission, error) {
	return authorizeOrgWriteResource(ctx, rt, oid)
}

// authorizeOrgWriteResource authorizes the given org to write the resources of the given type.
// NOTE: this is pretty much the same as AuthorizeWrite, in the case that the resource ID is ignored.
// Use it in the case that you do not know shich resource in particular you want to give access to.
func authorizeOrgWriteResource(ctx context.Context, rt manta.ResourceType, oid manta.ID) (manta.Authorizer, manta.Permission, error) {
	return authorize(ctx, manta.WriteAction, rt, nil, &oid)
}
