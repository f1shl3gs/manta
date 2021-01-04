package kv

import "context"

func Initial(ctx context.Context, kv SchemaStore) error {
	for _, bucket := range [][]byte{
		organizationBucket,
		organizationNameIndexBucket,
		nodeBucket,
		nodeOrgIndexBucket,
		nodeAddressIndexBucket,
		checkBucket,
		checkOrgIndexBucket,
		datasourceBucket,
		datasourceOrgIndexBucket,
		datasourceNameIndexBucket,
		notificationEndpointBucket,
		notificationEndpointNameIndexBucket,
		authorizationBucket,
		authorizationTokenIndexBucket,
		taskBucket,
		taskOrgIndexBucket,
		scraperTargetBucket,
		scraperTargetOrgIDBucket,
		userBucket,
		userNameIndexBucket,
		templateBucket,
		templateNameIndexBucket,
		datasourceBucket,
		datasourceNameIndexBucket,
		otclBucket,
		otclOrgIndex,
		dashboardBucket,
		dashboardOrgIndexBucket,
	} {
		err := kv.CreateBucket(ctx, bucket)
		if err != nil {
			return err
		}
	}

	/*err := kv.Update(ctx, func(tx Tx) error {
		b, err := tx.Bucket(dashboardBucket)
		if err != nil {
			return err
		}

		c, err := b.ForwardCursor(nil)
		if err != nil {
			return err
		}

		err = WalkCursor(ctx, c, func(k, v []byte) error {
			return b.Delete(k)
		})

		if err != nil {
			return err
		}

		b, err = tx.Bucket(dashboardOrgIndexBucket)
		if err != nil {
			return err
		}

		c, err = b.ForwardCursor(nil)
		if err != nil {
			return err
		}

		return WalkCursor(ctx, c, func(k, v []byte) error {
			return b.Delete(k)
		})
	})

	if err != nil {
		panic(err)
	}*/

	return nil
}
