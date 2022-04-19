#!/bin/bash

ip=$(ip address | grep -oP -m 1 '(?<=inet )(?!10|127)(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?=\/\d* brd \d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3} scope \w+ (?=eth\d+))')

echo "Swarm Init..."
docker swarm leave --force
docker swarm init --listen-addr $ip:2377 --advertise-addr $ip:2377
docker swarm join-token --quiet worker > /tmp/worker_token.txt
echo $ip > /tmp/manager_ip.txt

apt-get install s3cmd -y

# DigitalOcean droplets run updates after coming online, so we have to wait for it to finish to use apt-get
while [ $? != 0 ]
do
    sleep 10
    apt-get install s3cmd -y
done

echo "[default]
access_key = $ACCESS_KEY
access_token = 
add_encoding_exts = 
add_headers = 
bucket_location = fra1
ca_certs_file = 
cache_file = 
check_ssl_certificate = True
check_ssl_hostname = True
cloudfront_host = cloudfront.amazonaws.com
default_mime_type = binary/octet-stream
delay_updates = False
delete_after = False
delete_after_fetch = False
delete_removed = False
dry_run = False
enable_multipart = True
encoding = UTF-8
encrypt = False
expiry_date = 
expiry_days = 
expiry_prefix = 
follow_symlinks = False
force = False
get_continue = False
gpg_command = /usr/bin/gpg
gpg_decrypt = %(gpg_command)s -d --verbose --no-use-agent --batch --yes --passphrase-fd %(passphrase_fd)s -o %(output_file)s %(input_file)s
gpg_encrypt = %(gpg_command)s -c --verbose --no-use-agent --batch --yes --passphrase-fd %(passphrase_fd)s -o %(output_file)s %(input_file)s
gpg_passphrase = 
guess_mime_type = True
host_base = fra1.digitaloceanspaces.com
host_bucket = $BUCKET_NAME
human_readable_sizes = False
invalidate_default_index_on_cf = False
invalidate_default_index_root_on_cf = True
invalidate_on_cf = False
kms_key = 
limit = -1
limitrate = 0
list_md5 = False
log_target_prefix = 
long_listing = False
max_delete = -1
mime_type = 
multipart_chunk_size_mb = 15
multipart_max_chunks = 10000
preserve_attrs = True
progress_meter = True
proxy_host = 
proxy_port = 0
put_continue = False
recursive = False
recv_chunk = 65536
reduced_redundancy = False
requester_pays = False
restore_days = 1
restore_priority = Standard
secret_key = $SECRET_KEY
send_chunk = 65536
server_side_encryption = False
signature_v2 = False
signurl_use_https = False
simpledb_host = sdb.amazonaws.com
skip_existing = False
socket_timeout = 300
stats = False
stop_on_error = False
storage_class = 
urlencoding_mode = normal
use_http_expect = False
use_https = True
use_mime_magic = True
verbosity = WARNING
website_endpoint = http://%(bucket)s.s3-website-%(location)s.amazonaws.com/
website_error = 
website_index = index.html" > /root/.s3cfg

# manipulating persistent storage, shared with workers
s3cmd del s3://$BUCKET_NAME/worker_token.txt
s3cmd del s3://$BUCKET_NAME/manager_ip.txt
s3cmd put /tmp/worker_token.txt s3://$BUCKET_NAME
s3cmd put /tmp/manager_ip.txt s3://$BUCKET_NAME

echo "Sleeping for 30 seconds to let workers join swarm"
sleep 10 # wait for workers to join the swarm before deploying stack
echo "20 seconds remaining..."
sleep 10
echo "10 seconds remaining..."
sleep 5
echo "5 seconds remaining..."
sleep 5
docker stack deploy --compose-file docker-compose.yml $STACK_NAME
# services=$(docker service ls | grep -o -E minitwit_[a-zA-Z0-9_]+)
# for s in $services
# do
#     docker service scale $s=2
# done
