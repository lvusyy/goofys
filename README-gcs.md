# Google Cloud Storage (GCS)

Google Cloud Storage supports HTTP multi-range requests for improved performance with sparse file access patterns. Use the `--enable-multi-range` flag to enable this feature.

## Prerequisite

Service Account credentials or user authentication. Ensure that either the service account or user has the proper permissions to the Bucket / Object under GCS.

To have a successful mount, we require users to have object listing (`storage.objects.list`) permission to the bucket.

### Service Account credentials

Create a service account credentials (https://cloud.google.com/iam/docs/creating-managing-service-accounts) and generate the JSON credentials file.

### User Authentication and `gcloud` Default Authentication
User can authenticate to gcloud's default environment by first installing cloud sdk (https://cloud.google.com/sdk/) and running `gcloud auth application-default login` command.


## Using Goofys for GCS

### With service account credentials file
```
GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" goofys gs://[BUCKET] /path/to/mount
```

### With user authentication (`gcloud auth application-default login`)

```
goofys gs://[BUCKET] [MOUNT DIRECTORY]
```

### With multi-range requests enabled

```
GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" goofys --enable-multi-range gs://[BUCKET] /path/to/mount
```