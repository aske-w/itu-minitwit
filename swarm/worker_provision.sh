#!/bin/bash

echo "Opening firewall ports"
ufw allow 2376/tcp && \
ufw allow 2377/tcp && \
ufw allow 7946/tcp && \
ufw allow 7946/udp && \
ufw allow 4789/udp && \
ufw reload && \
ufw --force enable
echo "Finished opening firewall ports"
#!/bin/bash
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

s3cmd --force get s3://$BUCKET_NAME/worker_token.txt /root/
s3cmd --force get s3://$BUCKET_NAME/manager_ip.txt /root/

docker swarm leave
docker swarm join --token $(cat /root/worker_token.txt) $(cat /root/manager_ip.txt):2377