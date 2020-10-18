# s3-restore

A tiny utility to restore deleted objects from a S3 version enabled bucket.

## Background

When an object is deleted from a version-enabled bucket, S3 creates a delete marker for the object.
The delete marker becomes the current version of the object, and the actual object becomes the
previous version. The deleted objects from a version enabled bucket can be restored in by either:

1. Copying the previous version of the object to the same or another bucket.
2. Remove the delete marker, this requires running the script as the user who is the AWS account
owner or the user who created the bucket.

## Install