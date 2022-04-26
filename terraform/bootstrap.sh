#!/bin/bash

echo -e "\n--> Bootstrapping Minitwit\n"

echo -e "\n--> Loading environment variables from secrets file\n"

source secrets

echo -e "\n--> Checking that environment variables are set\n"
# check that all variables are set
[ -z "$TF_VAR_DOTOKEN" ] && echo "TF_VAR_DOTOKEN is not set" && exit
[ -z "$BUCKET_NAME" ] && echo "BUCKET_NAME is not set" && exit
[ -z "$STATE_FILE" ] && echo "STATE_FILE is not set" && exit
[ -z "$ACCESS_KEY" ] && echo "ACCESS_KEY is not set" && exit
[ -z "$SECRET_KEY" ] && echo "SECRET_KEY is not set" && exit

echo $ACCESS_KEY
echo $SECRET_KEY

echo -e "\n--> Initializing terraform\n"
# initialize terraform
terraform init \
    -backend-config "bucket=$BUCKET_NAME" \
    -backend-config "key=$STATE_FILE" \
    -backend-config "access_key=$ACCESS_KEY" \
    -backend-config "secret_key=$SECRET_KEY"

# check that everything looks good
echo -e "\n--> Validating terraform configuration\n"
terraform validate

# create infrastructure
echo -e "\n--> Creating Infrastructure\n"
terraform apply -auto-approve