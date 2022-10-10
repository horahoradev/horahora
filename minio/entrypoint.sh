#!/bin/bash
set -euo pipefail

# Run storage setup after minio starts (meaning, this runs in the background)
{
  mc alias set minio http://localhost:9000 $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD --api S3v4

  until mc du minio; do
    echo "Waiting for minio to start..."
    sleep 1
  done

  echo "Creating bucket"
  # TODO(ivan): How do we make this run once?
  mc mb minio/$BUCKET_NAME || true
  # make it so we can download files from it
  mc anonymous set download minio/$BUCKET_NAME
} &

exec minio server /data
