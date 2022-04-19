#! /bin/bash
# Provisioner script for Minitwit Docker swarm since Vagrant decided not to work

if [ ! -f `which jq` ]
then
  echo "Dependency jq is not installed. Exiting..."
  exit 1
fi

manager_tags='["swarm", "manager"]'
manager_count=1
manager_size="s-2vcpu-4gb"
manager_name="manager"

worker_tags='["swarm", "worker"]'
worker_count=1
worker_size="s-2vcpu-4gb"
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
echo "Fetching existing droplet ID's"
existing_droplets=$(curl -# -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  "https://api.digitalocean.com/v2/droplets?page=1&per_page=200"\
  | jq '.droplets | .[] | .id')
echo "Fetched existing droplet ID's"

# delete existing droplets
echo "Deleting existing droplets"
for id in $existing_droplets 
do
  echo "Removing droplet $id"
  curl -# -s -X DELETE \
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
echo "Deleted existing droplets"

# get ssh all key ids in team
echo "Fetching SSH keys for droplets"
ssh_ids=$(curl -# -X GET \
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

# create manager droplet and save its id
name="${manager_name}-1"
echo "Creating droplet $name"
manager_id=$(curl -# -S -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
-d "{\"name\":\"$name\",
  \"region\":\"$region\",
  \"size\":\"${manager_size}\",
  \"image\":\"${image}\",
  \"ssh_keys\":${ssh_ids},
  \"user_data\":\"$(cat authorized_keys > /root/.ssh/authorized_keys)\",
  \"tags\":${manager_tags}}" \
  "https://api.digitalocean.com/v2/droplets" \
  | jq -r '.droplet | .id')
if [ $? = 0 ]
then
  echo "Created droplet $name"
else
  echo "Something went wrong during API request to DigitalOcean. Exiting..."
  exit 1
fi


# for i in `seq -w 1 $manager_count`
# do
#   name="${manager_name}-${i}"
#   echo "Creating droplet $name"
#   curl -S -X POST \
#   -H "Content-Type: application/json" \
#   -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
#   -d "{\"name\":\"$name\",
#     \"region\":\"$region\",
#     \"size\":\"${manager_size}\",
#     \"image\":\"${image}\",
#     \"ssh_keys\":${ssh_ids},
#     \"user_data\":\"$(cat ./ufw_script.sh)\",
#     \"tags\":${manager_tags}}" \
#   "https://api.digitalocean.com/v2/droplets" \
#   1> /dev/null
#   if [ $? = 0 ]
#   then
#     echo "Created droplet $name"
#   else
#     echo "Something went wrong during API request to DigitalOcean. Exiting..."
#     exit 1
#   fi
# done

#exit 0;

# TODO: get first manager ip and SCP the worker token saved in /tmp/worker_token.txt to this machine
# or copy it to the swarm-storage space


# create worker droplets
for i in `seq -w 1 $worker_count`
do
  name="${worker_name}-${i}"
  echo "Creating droplet $name"
  curl -# -S -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  -d "{\"name\":\"$name\",
    \"region\":\"$region\",
    \"size\":\"${worker_size}\",
    \"image\":\"${image}\",
    \"ssh_keys\":${ssh_ids},
    \"user_data\":\"$(cat authorized_keys > /root/.ssh/authorized_keys)\",
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

# get floating IP
echo "Fetching floating IP"
floating_ip=$(curl -# -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  "https://api.digitalocean.com/v2/floating_ips?page=1&per_page=20" \
  | jq -r '.floating_ips | .[0] | .ip')
if [ $? = 0 ]
then
  echo "Fetched floating IP $floating_ip"
else
  echo "Something went wrong during API request to DigitalOcean. Exiting..."
  exit 1
fi

echo "Assigning floating IP $floating_ip to manager droplet. This will usually take 1-3 minutes."
echo "Assigning floating IP to manager droplet in 30 seconds"
sleep 5
echo "Assigning floating IP to manager droplet in 25 seconds"
sleep 5
echo "Assigning floating IP to manager droplet in 20 seconds"
sleep 5
echo "Assigning floating IP to manager droplet in 15 seconds"
sleep 5
echo "Assigning floating IP to manager droplet in 10 seconds"
sleep 5
echo "Assigning floating IP to manager droplet in  5 seconds"
sleep 5

# assign floating IP to manager
echo "Attempting to assign floating IP to $manager_name"
status_code=$(curl -o /dev/null -w "%{http_code}" -# -S -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
  -d "{\"type\":\"assign\",\"droplet_id\":$manager_id}" \
  "https://api.digitalocean.com/v2/floating_ips/$floating_ip/actions")

echo "Attempt returned status code $status_code"
if [[ $status_code = "422" ]]
then
  echo "The script will retry while the status code is 422"
fi

while [[ $status_code = "422" ]]
do
  echo "Reattempting assigning floating IP to manager droplet in 15 seconds"
  sleep 5
  echo "Reattempting assigning floating IP to manager droplet in 10 seconds"
  sleep 5
  echo "Reattempting assigning floating IP to manager droplet in  5 seconds"
  sleep 5
  echo "Attempting to assign floating IP to $manager_name"
  status_code=$(curl -o /dev/null -w "%{http_code}" -# -S -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" \
    -d "{\"type\":\"assign\",\"droplet_id\":$manager_id}" \
    "https://api.digitalocean.com/v2/floating_ips/$floating_ip/actions")
  echo "Returned $status_code"
done

