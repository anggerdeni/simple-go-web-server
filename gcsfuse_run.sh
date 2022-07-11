#!/usr/bin/env sh
set -eo pipefail

# Create mount directory for service
mkdir -p $MNT_DIR

echo "Mounting GCS Fuse."
gcsfuse --debug_gcs --debug_fuse $BUCKET $MNT_DIR 
echo "Mounting completed."

# Run the web service on container startup
exec /bin/server &

# Exit immediately when one of the background processes terminate.
wait -n