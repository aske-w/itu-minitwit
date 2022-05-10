resource "digitalocean_floating_ip" "public-ip" {
  region = var.region
}

resource "digitalocean_floating_ip_assignment" "public-ip" {
  ip_address = digitalocean_floating_ip.public-ip.ip_address
  droplet_id = digitalocean_droplet.swarm-manager.id
}

output "public_ip" {
  value = digitalocean_floating_ip.public-ip.ip_address
}