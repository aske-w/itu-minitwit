# MacOS

## Vagrant

- `brew install virtualbox`
- `brew install virtualbox-extension-pack`
- `brew install vagrant`
- `brew install vagrant-completion`

If no SSH key, generate one with `ssh-keygen -t rsa`

Check if installed correctly with `vagrant --version`

Documentation: https://www.vagrantup.com/docs

### Plugins

- `vagrant plugin install vagrant-digitalocean`
- `vagrant plugin install vagrant-scp`
- `vagrant plugin install vagrant-vbguest`
- `vagrant plugin install vagrant-reload`
- WMWARE: `vagrant plugin install vagrant-vmware-desktop`

## VirtualBox networking

https://discuss.hashicorp.com/t/vagrant-2-2-18-osx-11-6-cannot-create-private-network/30984/22

Create the file `/etc/vbox/networks.conf` if it doesnt exist
Add `* 0.0.0.0/0 ::/0` to the file and save.

Now we can create networks in the Vagrant file:

```ruby

 #...
 config.vm.network "private_network", type: "dhcp" # <---
 #...
  config.vm.define "webserver", primary: true do |server|
        server.vm.network "private_network", ip: "192.168.62.1" # <---
        server.vm.network "forwarded_port", guest: 8080, host: 8080

```
