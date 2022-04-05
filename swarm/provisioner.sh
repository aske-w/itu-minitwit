#! /bin/bash
# Provisioner script for Minitwit Docker swarm since Vagrant decided not to work

if [ ! -f `which jq` ]
then
  echo "Dependency jq is not installed. Exiting..."
  exit 1
fi

manager_tags='["swarm", "manager"]'
manager_count=1
manager_size="s-1vcpu-1gb"
manager_name="manager"

worker_tags='["swarm", "worker"]'
worker_count=2
worker_size="s-1vcpu-1gb"
worker_name="worker"

image="docker-18-04"
region="fra1"

# while getopts 't:' OPTION; do
#     case "$OPTION" in
#         t)
#             ;;
#         *)
#             echo "Usage: $0 [-t digital_ocean_token]" >&2
#             exit 1
#         ;;
#     esac
# done

DIGITALOCEAN_TOKEN=$1 # provided as argument

# get existing droplet ids
existing_droplets=$(curl -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  "https://api.digitalocean.com/v2/droplets?page=1&per_page=200"\
  | jq '.droplets | .[] | .id')

# delete existing droplets
for id in $existing_droplets 
do
  echo "Removing droplet $id"
  curl -s -X DELETE \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
    "https://api.digitalocean.com/v2/droplets/$id"
  if [ $? = 0 ]
  then
    echo "Removed droplet $id"
  else
    echo "Something went wrong during API request to DigitalOcean. Exiting..."
    exit 1
  fi
done

# get ssh all key ids in team
echo "Fetching SSH keys for droplets"
ssh_ids=$(curl -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  "https://api.digitalocean.com/v2/account/keys" \
  | jq '[.ssh_keys|.[]|.id'])
if [ $? != 0 ]
then
  echo "Error fetching SSH keys for droplets. Exiting..."
  exit 1
fi
echo "Fetched SSH keys for droplets"

# create manager droplets
for i in `seq -w 1 $manager_count`
do
  name="${manager_name}-${i}"
  echo "Creating droplet $name"
  curl -S -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  -d "{\"name\":\"$name\",
    \"region\":\"$region\",
    \"size\":\"${manager_size}\",
    \"image\":\"${image}\",
    \"ssh_keys\":${ssh_ids},
    \"user_data\":null,
    \"tags\":${manager_tags}}" \
  "https://api.digitalocean.com/v2/droplets" \
  1> /dev/null
  if [ $? = 0 ]
  then
    echo "Created droplet $name"
  else
    echo "Something went wrong during API request to DigitalOcean. Exiting..."
    exit 1
  fi
done

# create worker droplets
for i in `seq -w 1 $worker_count`
do
  name="${worker_name}-${i}"
  echo "Creating droplet $name"
  curl -S -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  -d "{\"name\":\"$name\",
    \"region\":\"$region\",
    \"size\":\"${worker_size}\",
    \"image\":\"${image}\",
    \"ssh_keys\":${ssh_ids},
    \"user_data\":null,
    \"tags\":${worker_tags}}" \
  "https://api.digitalocean.com/v2/droplets" \
  1> /dev/null
  if [ $? = 0 ]
  then
    echo "Created droplet $name"
  else
    echo "Something went wrong during API request to DigitalOcean. Exiting..."
    exit 1
  fi
done