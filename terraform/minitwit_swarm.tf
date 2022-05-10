
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
  # provisioner "file" {
  #   source = "./authorized_keys"
  #   destination = "/root/.ssh/authorized_keys"
  # }

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
    command = "mkdir temp"
  }
  provisioner "local-exec" {
    command = "chmod -R 600 ssh_key"
  }
  provisioner "local-exec" {
    command = "ssh -o 'StrictHostKeyChecking no' root@${self.ipv4_address} -i ssh_key/terraform 'docker swarm join-token worker -q' > temp/worker_token.txt"
  }
  provisioner "local-exec" {
    command = "ssh -o 'StrictHostKeyChecking no' root@${self.ipv4_address} -i ssh_key/terraform 'echo \"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDFEH+GGxAlOd+HcZWfJs77G0gC7O52OZXn3+WaGt+YJ+krQazfW4Icy1SoIssnII0ZMXHNXrzbiAc/EIZueTCPyqzLyjc3PpBXcUcNYaE3qlJibOH2oF3HJc0KzaVAyN4XG6oDvgymX85UC63hbPEwnd7N6ju++NvbkovS/DoSeyTcI24vbdh8vL0dr7SjI9e71YMtmsfj5BI+6KaoJX3RE+gAg+GiGdLFr17vUPWCjiXtsb1qD/UOxTPTlRkP7aEVIYlxQjKzj3XicEoIZF7IttB654WbE0N7ogmcQty/PaiNYlAEcSzXVNRniv7X9ZLNpKZSzhp0DRr2qbf4e9VnWRQJ9qBMKfWzKsv+0zHFNYK7PUnXbZvjg57+M2gUAMTqkQhs4qiDK4WtYtyQmPT5Ca4ywFaAicO0WHWbeVbvpLmv+QslIsw0qQVKfYpOJvCK8evkCnH/Ufh07hy51owfZn/0lWTx2L8OAmq3fdrN5wO0j6m+mfsk/69SzMj2Lns= wachs@WachsMac.local\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC6OwiKMLxGtc48mCHZ0wJ84OqKmISEfs9HIPqoslSJLlrqjJjD9KZCE2B8MiV2Luvd76XpECK0rUJoIl9dY9m+hzDJ8nXK64/FsDzihQyPrhIwL2rM1veQz6lISQnrAGHpXsNWgt326YGlBh0N2OJLRa3AmzJxIwOjOYf8xq5PKUdb2S45zk9wDWN/0PPAkl0Ee7RRyjaMW1IyulkosJk7CYaP2049aznJCvd5qmrxdEu1nXJUEY8gXzSHqRNEkKw51j6X1/DtGXJ8sbBp+KJEjVcmyn92y/yRPO/1kDwJwADlTvXqgWBp3+317Iljj3W83J1czmbhK5pFkh5/W2YHMKukquMpk79sL/2LatZAit8TJaCTdncC7DA5cH2AowPzps2EW81vOlp3trqjf95A2aZhiYvk3tJq3TWmL75XuQN/t9dXDjerf5k6ETNL5X8zCpkjRby7xGH3Pk0uIY18JCr8C2ejFue/5UiSxcYuWkOERIIGmzIH8cX9hG9gKPM= askew@LAPTOP-17F8SRCT\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCtXVIGA58bpcw+ubS3Zgc5w5H40v/o1wYZyGYsIA5R0A5dgz4hZXqTyY0Zp9BSkWOTIQGvIy+tCQM5qttJ2HcFbsyh0mvZ8UmoDbhG/DYu80y1R4NAnILkmLU5eyqIZ+Qfc4YCyJDXLqbtQ6NiQo4suw2pmHwlXihKyNhcPTKFUJp5MbIcgqcgebg+sKJo95vcp4RcvVH/XL+qB671YFEUiXqwu1Ygpc8Vv/JJ+gpaAKNfaY8HW3HWhdlpH224VYVEpOeqCKRit6Au9/aML4t36qydqY+Om0xgwWzdTFMp3pnIL6TjB8QHKjsDnlg+MK5x+x/osocKeE+FkOxOED6x6hMBPT7nv5fcU5cLn8xrRHras2iU0tSQpZGH8kp4p3bz0gcPDUC7m/Zv+htkWab4H2BSMWBmkOeZ7MYFceXpyczD+QgWQbpJVIklVadhvaJwwZvvoGf8OQ9oU0KB5mHNB9rzzXu1n2xgezhVdNCEM7B247pKcKJ8z8P34q3GTUs= wermuth@Alexanders-MacBook-Pro.local\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDZG1CcrtBVuw0dpiQaNj80HHrB2dinEcuvkr0BP9qIKtHwzkQsKdYu6OTG5h+qk+dWBf9Dn3B8BxPVyC0Hi1HdXr70uuaumSvy/nqkhI+KSxEaHqyWIWTl4cpj3XSrg9++i8Qd4Ic0dPMmG+e23m1GCC2r6EsoB8w3tutTQzPegXkxxvyR7yl3jWU2Pu4/3jcnLfT4P2ot7I9WgxQcs8/y3IKiAnzjJxV5Rf3UFRVwebz3i1+M+UVL8JLE6ZfxMQh/sGFSmS+N41w7tq3Hi4vNrlzLt+p18v6LH/A0f6Q4iPi26D0F4GnzamVX4KJG1wX8cffpY4nfYUbl1Ph1V5ynTPp/J6lIAJDVYLDDzJcrYWFKNiV3cqVuaLi5LHtzt8PsPc4IznEHUgUgF3gvS+KROcE8EnxMDSl/ZPPGgriUTC5UB+FCCjjm3DptKg5Y8MjxRaw0JsE4H144NiF66Jtjy+w2cqvyVqZrfoNP5AJvdOiaVV8u9u2BrbpDIJ/Br8U= wachs@WachsMac.local\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCayWaNTtF0r9W1Lbza7239QHa/8fTdjs/U4ZIDq+pY8d7eilZv0eVuGjVw/nE8lNE30YOh9Y2Abh3E7B7HLoyop9kkPv0nEhgbd6cUwSEetNnxA+6YxSXeKe35vRed+JddX5z1tY/BaffsC4ruHu3gi2XgdPwKBMSwt7Vtc17dRFpblTy2RbJ6RWPt+uXjqACal3bKNBN0hkx6+Br4p2qdwdD4tsVs/YhihutWW6P/KeFgWQ8TnYdR3cuj1s+eIhWMfyIJl2FkMrv+nFvAO6mGaG+xNj/royGVG43JwmVvY9KBIYXUwsfZi9CxhmiKkRlJXmvhdywmBUXs/P4HSkGVKNYcanB8zKUSzbV0gPGMMso4y/wTR4tC/bbTliPKOu1rKo6ONpInW8+Qwg51xeFk3PaUMM19ESGb8d8/JiQT3HZRu9KplUH+qHlyC/y8E6ugoFFqVQEZyv2TZzQ6dgP3Ns0cdpeI9e/3LZuRqFLMF8xxVLfEJMUf1Gz7INXfLik= christian@feedr.com\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDssQaeAVaIndQdY8jieLGa7hqcvfIcgtQmdZJsEWY82mQ/LiJr75c2LEqfTduKJ+xEToKRSgTJx/dWADs2I+UfU92yyTIR0tV1H+r1Pt/QEb4Zf2h+uU18X35w+rY648yim21/LElHOXMHFi1I24cbv1kpJv+/+A4PhAOM8ECRY44Ro3s7kURAcc5BbDRuX+2Mk53tAojYZJIhqnkxSpE/62DOF2rEG7Mq07ToAEjL77ta41E52Si3h1WxtnwbMDC4x3ZYAo7NOd5FWomA15zx4j9aH45WHS06dq95Wvsd4c1qR0WCVxSSXWF2D214M1886MZEEgKjjD8llLqkEXgoPn0Qzuma7afk0gOGOa+NX3/C0PEmO2ZZnVdojs7KE9mAae2ErGbO3Bu38OD184QjUWKv+BtwBDzbJNDr2wVEhHgccy7C1olvVTmmfqC9/aFOosSOIwEKi32fkNZy8AZEJlDrFbirmpeM+8h/5S1wNn7NZcQTAaOuokWG/UUVRYU= jacobmolby@Moelbys-MacBook.local\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDwdyEj/oYh5wNmnFHei+jMB6eNPBOs60JXoRCB/A0eB3Tmf5SXXnc3ILeABjjrrv5Y5MfRnnKRu/hj9ttjGm+AK5P+GPAWSc8QIPy8ZlJu9TYHLgf3cQBCaOVbUMCx/tMB43bIT1/Ari2VRbt2ZzuZkucPZx8bJyibzor0UUJc+QbZw37BoDlDZWjfZ6RwdukC+Do0WL3Cf7lx+uc5NIEWp7aKU9F+FGRVrSrWLCdwYU6RfsDCYMXLRjKxnUDGDbQTYRvTIlgI63n5l2OGoy62hsG5S7Kg2SdeGUyTLH2SCSmbu9CT32t4TTaaPSnuVafFc9j0d6qwJS9YI2jHJ9ZNtOk+9OERTNM3tGaX0E2cJLGVj1C57sdwDVVXykkIcpour+I6z+PxjBu2sPj9hLqKo6kC8ggiaRz+dbAfmxt+USkrj5NJgcq+KsgcOpo7USGbH2h4n88+ZsMuzhpbsAdV+kCDEMcvlZ96ZLGgoiMk5tctGjN9GAd9c7gfhfkv5nP72svOg/ASeNuh3x9k0D8ivmeJTLXFqnToeXYz30NgXRSpws57zsPvZQN6a31tunrVRgkcm+zVu6mGX9k2h1LuEfH/Z8bGapLVwjkpK6Vmcq5MYauXW30WKL9tYlkZlrwU0i7icAY80UDltKnMc4Pc++/0F2dDPWqwRgDRDxMAmQ== lars@LAPTOP-PNLOF9HB\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDBD+fJHYMRouycCSkRZAULVz5puO7/YOO21HEHF4jnriZIbKruPCHM5xpjCU3o4WD2Zpye6DtuNwhbcLeCZsbU/WY+IPeHUV2kLeGfU+jKFS8S6jW0g6HUKcjSlP8wyI/OCooecU+QdaoOHx9uoCGGxps8f9KtVJC0rd0G0zmjoDPd2VAlFaAa3QcKNeuaPFtcDv1mJifz6X13zERFe3WM5DGpGH0E9Ii3y1g1l9vz5xbt/kCji5bqg8d+6nvGk5JHCfOpDLjxd5LUxgnusthO8GZkCHcOcTjrG1Fz2iWJvDTOxd8BUdAzW1qeL8wGrCjKsv2oBlFvldCvzB8PPYyoHiBSXQh3h/znUr6O0wRmcqcJ6ntb1PUAagF6XfBqqEBno2kjk+hWuMSlWy2vxojBo4+xCYISaeCzstxDCFxLVCxrXtnefP4PP1KZEI5ltmGaqEBzv1ig2xJv6aTkFAGYV1I4xbwpEI3fa0RY49r1nDLj8S1OqPJv9/71zhY0hDzjNVfxlwe8sR8iCFXY0kKJ0EUINNSf4284tCnVhZBBDhm8vM7IFtl69tIrDFF8i2rWB3qOkB2IQUpWMTDmtIVS40G7EN1TO3RNEd4hlMv96DNgB/ly2QZUm2UyKij8OR0J4tgBcOJCKR4ZZ6oSgV0/WFHQ7+IXJ0DIxx4BxBGFSw== tobia@LAPTOP-G5JPLPAG\" >> /root/.ssh/authorized_keys'"
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