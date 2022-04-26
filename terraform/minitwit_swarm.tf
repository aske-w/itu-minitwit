
#  _ __ ___   __ _ _ __   __ _  __ _  ___ _ __
# | '_ ` _ \ / _` | '_ \ / _` |/ _` |/ _ \ '__|
# | | | | | | (_| | | | | (_| | (_| |  __/ |
# |_| |_| |_|\__,_|_| |_|\__,_|\__, |\___|_|
#    
resource "digitalocean_droplet" "swarm-manager" {
  image = var.image
  name = var.manager_name
  region = var.region
  size = var.manager_size
  # public_key = var.pub_key

  connection {
    user = "root"
    host = self.ipv4_address
    type = "ssh"
    private_key = var.pvt_key
    timeout = "1m"
  }

  # Prometheus
  provisioner "file" {
    source = "../swarm/prometheus"
    destination = "/root"
  }

  # Grafana
  provisioner "file" {
    source = "../swarm/grafana"
    destination = "/root"
  }

  # Add file beat
  provisioner "file" {
    source = "../filebeat/filebeat.yml"
    destination = "/root/filebeat.yml"
  }

  provisioner "file" {
    source = "../swarm/docker-compose.yml"
    destination = "/root/docker-compose.yml"
  }

  provisioner "file" {
    source = "./.s3cfg"
    destination = "/root/.s3cfg"
  }

  provisioner "remote-exec" {
    inline = [
      "echo \"Opening firewall ports\"",
      "ufw allow 2376/tcp",
      "ufw allow 2377/tcp",
      "ufw allow 7946/tcp",
      "ufw allow 7946/udp",
      "ufw allow 4789/udp",
      "ufw reload",
      "ufw --force enable",
      "echo \"Finished opening firewall ports\"",

      "ip=$(ip address | grep -oP -m 1 '(?<=inet )(?!10|127)(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})(?=\\/\\d* brd \\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3} scope \\w+ (?=eth\\d+))')",

      "echo \"Swarm Init...\"",
      "docker swarm leave --force",
      "docker swarm init --listen-addr $ip:2377 --advertise-addr $ip:2377",
      "docker swarm join-token --quiet worker > /tmp/worker_token.txt",
      "echo $ip > /tmp/manager_ip.txt",

      "apt-get install s3cmd -y",

      # DigitalOcean droplets run updates after coming online, so we have to wait for it to finish to use apt-get
      "while [ $? != 0 ];do;sleep 10;apt-get install s3cmd -y;done;",

      # manipulating persistent storage, shared with workers
      "s3cmd del s3://${BUCKET_NAME}/worker_token.txt",
      "s3cmd del s3://${BUCKET_NAME}/manager_ip.txt",
      "s3cmd put /tmp/worker_token.txt s3://${BUCKET_NAME}",
      "s3cmd put /tmp/manager_ip.txt s3://${BUCKET_NAME}",

      "echo \"Sleeping for 30 seconds to let workers join swarm\"",
      "sleep 10 # wait for workers to join the swarm before deploying stack",
      "echo \"20 seconds remaining...\"",
      "sleep 10",
      "echo \"10 seconds remaining...\"",
      "sleep 5",
      "echo \"5 seconds remaining...\"",
      "sleep 5",
      "docker-compose pull",
      "docker stack deploy --compose-file docker-compose.yml ${STACK_NAME}"
    ]
  }
}

#                     _
# __      _____  _ __| | _____ _ __
# \ \ /\ / / _ \| '__| |/ / _ \ '__|
#  \ V  V / (_) | |  |   <  __/ |
#   \_/\_/ \___/|_|  |_|\_\___|_|
#

resource "digitalocean_droplet" "swarm-worker" {
  # create workers after the leader
  depends_on = [digitalocean_droplet.swarm-manager]

  # number of vms to create
  count = 1

  image = var.image
  name = "${var.worker_name}-${count.index}"
  region = var.region
  size = var.worker_size
  # add public ssh key so we can access the machine
  # ssh_keys = [digitalocean_ssh_key.minitwit.fingerprint]

  # specify a ssh connection
  connection {
    user = "root"
    host = self.ipv4_address
    type = "ssh"
    private_key = var.pvt_key
    timeout = "2m"
  }

  # Prometheus
  provisioner "file" {
    source = "../swarm/prometheus"
    destination = "/root"
  }

  # Grafana
  provisioner "file" {
    source = "../swarm/grafana"
    destination = "/root"
  }

  # Add file beat
  provisioner "file" {
    source = "../filebeat/filebeat.yml"
    destination = "/root/filebeat.yml"
  }

    provisioner "file" {
    source = "./.s3cfg"
    destination = "/root/.s3cfg"
  }
  
  provisioner "remote-exec" {
    inline = [
      "echo \"Opening firewall ports\"",
      "mkdir /var/prom_data",
      "chown nobody /var/prom_data",
      
      "chown root filebeat.yml",
      "chmod go-w filebeat.yml",

      "ufw allow 2376/tcp",
      "ufw allow 2377/tcp",
      "ufw allow 7946/tcp",
      "ufw allow 7946/udp",
      "ufw allow 4789/udp",
      "ufw reload",
      "ufw --force enable",
      "echo \"Finished opening firewall ports\"",
      
      "apt-get install s3cmd -y",
      # DigitalOcean droplets run updates after coming online, so we have to wait for it to finish to use apt-get
      "while [ $? != 0 ];do;sleep 10;apt-get install s3cmd -y;done;",

      "s3cmd --force get s3://${BUCKET_NAME}/worker_token.txt /root/",
      "s3cmd --force get s3://${BUCKET_NAME}/manager_ip.txt /root/",

      "docker swarm leave",
      "docker swarm join --token $(cat /root/worker_token.txt) $(cat /root/manager_ip.txt):2377"
    ]
  }
}

output "swarm-manager-ip-address" {
  value = digitalocean_droplet.swarm-manager.ipv4_address
}

output "swarm-worker-ip-address" {
  value = digitalocean_droplet.swarm-worker.*.ipv4_address
}