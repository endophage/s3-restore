package restore

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func createS3Client(region string) *s3.S3 {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	return s3.New(sess)
}

// ObjsUingCopy - Restores objects by copying previous versions of the objects.
func ObjsUingCopy(bucket string, prefix string, region string, dstBucket string, dstPrefix string, dryrun bool) {
	log.Infof("Restoring objects by copying previous versions of the objects.")

	svc := createS3Client(region)
	vers, dels := getObjectVersions(svc, bucket, prefix)
	objs := make([]s3.ObjectVersion, 0)

	for _, ver := range *vers {
		for _, del := range *dels {
			if *ver.Key == *del.Key {
				findOrInsertObj(&objs, &ver)
			}
		}
	}

	for _, obj := range objs {
		copyVersionedObj(svc, bucket, *obj.Key, *obj.VersionId, dstBucket, dstPrefix, dryrun)
	}
}

// ObjsUsingDel - Restores objects by removing the DELETE markers.
func ObjsUsingDel(bucket string, prefix string, region string, dryrun bool) {
	log.Infof("Restoring objects by removing DELETE markers.")

	svc := createS3Client(region)
	_, objs := getObjectVersions(svc, bucket, prefix)

	for _, obj := range *objs {
		deleteVersionedObj(svc, bucket, *obj.Key, *obj.VersionId, dryrun)
	}
}

func getObjectVersions(svc *s3.S3, bucket string, prefix string) (*[]s3.ObjectVersion, *[]s3.DeleteMarkerEntry) {
	log.WithField("bucket", bucket).Infof("Listing object versions in bucket")

	vers := make([]s3.ObjectVersion, 0)
	dels := make([]s3.DeleteMarkerEntry, 0)

	params := s3.ListObjectVersionsInput{}
	params.SetBucket(bucket)
	if len(prefix) > 0 {
		params.SetPrefix(prefix)
	}

	svc.ListObjectVersionsPages(&params, func(res *s3.ListObjectVersionsOutput, lastPage bool) bool {
		for _, version := range res.Versions {
			if !*version.IsLatest {
				vers = append(vers, *version)
			}
		}

		for _, marker := range res.DeleteMarkers {
			if *marker.IsLatest {
				dels = append(dels, *marker)
			}
		}

		return true
	})

	if len(vers) > 0 {
		log.WithField("count", len(vers)).Infof("Found objects with previous versions")
	} else {
		log.Infof("There are no objects with previous versions")
	}

	if len(dels) > 0 {
		log.WithField("count", len(dels)).Infof("Found objects with Delete Markers")
	} else {
		log.Infof("There are no objects with Delete Markers")
	}

	return &vers, &dels
}

func copyVersionedObj(svc *s3.S3, srcBucket string, key string, vid string, dstBucket string, dstPrefix string, dryrun bool) {
	if len(dstBucket) == 0 {
		dstBucket = srcBucket
	}
	if len(dstPrefix) > 0 {
		key = fmt.Sprintf("%s/%s", dstPrefix, key)
	}

	copySource := fmt.Sprintf("%s/%s?versionId=%s", srcBucket, key, vid)

	log.WithFields(log.Fields{
		"src_obj": copySource,
		"dst_obj": fmt.Sprintf("%s/%s", dstBucket, key),
	}).Infof("Attempting to copy object.")

	if !dryrun {
		params := s3.CopyObjectInput{}
		params.SetCopySource(copySource)
		params.SetBucket(dstBucket)
		params.SetKey(key)

		_, e := svc.CopyObject(&params)
		if e != nil {
			log.Errorf("Failed copying the obj, reason: %s", e)
		}
	}
}

func deleteVersionedObj(svc *s3.S3, bucket string, key string, vid string, dryrun bool) {
	log.WithFields(log.Fields{"key": key, "version_id": vid}).Infof("Attempting to remove delete marker of the object")

	if !dryrun {
		params := s3.DeleteObjectInput{}
		params.SetBucket(bucket)
		params.SetKey(key)
		params.SetVersionId(vid)

		_, e := svc.DeleteObject(&params)
		if e != nil {
			log.Errorf("Failed deleting the obj, reason: %s", e)
		}
	}

}

func findOrInsertObj(objs *[]s3.ObjectVersion, obj *s3.ObjectVersion) {
	objExists := false

	for i, o := range *objs {
		if *obj.Key == *o.Key {
			objExists = true

			if (*obj.LastModified).After(*o.LastModified) {
				(*objs)[i] = *obj
			}
		}
	}

	if !objExists {
		*objs = append(*objs, *obj)
	}
}
