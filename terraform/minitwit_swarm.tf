
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
  ssh_keys = [digitalocean_ssh_key.minitwit.fingerprint]
  # public_key = var.pub_key

  connection {
    user = "root"
    host = self.ipv4_address
    type = "ssh"
    private_key = file(var.pvt_key)
    timeout = "5m"
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

      "ip=${digitalocean_droplet.swarm-manager.ipv4_address}",

      "echo \"Swarm Init...\"",
      "docker swarm leave --force",
      "docker swarm init --listen-addr $ip:2377 --advertise-addr $ip:2377",
    ]
  }

  # save the worker join token
  provisioner "local-exec" {
    command = "ssh -o 'StrictHostKeyChecking no' root@${self.ipv4_address} -i ssh_key/terraform 'docker swarm join-token worker -q' > temp/worker_token.txt"
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
  ssh_keys = [digitalocean_ssh_key.minitwit.fingerprint]

  # specify a ssh connection
  connection {
    user = "root"
    host = self.ipv4_address
    type = "ssh"
    private_key = file(var.pvt_key)
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
    source = "./temp/worker_token.txt"
    destination = "/root/worker_token.txt"
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

      "docker swarm leave",
      "docker swarm join --token $(cat /root/worker_token.txt) ${digitalocean_droplet.swarm-manager.ipv4_address}:2377"
    ]
  }

  # TODO this will cause workers except the first to probably not have services running in them
  # provisioner "local-exec" {
  #   command = "ssh -o 'StrictHostKeyChecking no' root@${digitalocean_droplet.swarm-manager.ipv4_address} -i ssh_key/terraform 'docker stack deploy --compose-file /root/docker-compose.yml minitwit'"
  # }
}

output "swarm-manager-ip-address" {
  value = digitalocean_droplet.swarm-manager.ipv4_address
}

output "swarm-worker-ip-address" {
  value = digitalocean_droplet.swarm-worker.*.ipv4_address
}