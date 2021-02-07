#!/bin/bash

if [[ $1 != @(read|write) ]]; then
    echo "Deployment argument can be only [read] or [write]. got: [$1]"
fi

echo -e "Deploying function [$1] to Google Cloud Function." 
# gcloud functions deploy $1 --runtime go113 --trigger-http --allow-unauthenticated
