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

	return nil
}
