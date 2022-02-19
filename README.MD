# MacOS

## Vagrant

- `brew install virtualbox`
- `brew install virtualbox-extension-pack`
- `brew install vagrant`
- `brew install vagrant-completion`

If no SSH key, generate one with `ssh-keygen -t rsa`

Check if installed correctly with `vagrant --version`

Documentation: https://www.vagrantup.com/docs

After you have installed Vagrant and created your Vagrantfile, you can start it with `vagrant up`
To delete the vm again write `vagrant destroy`.
To run the provision scripts again, on an existing vm run `vagrant provision --provision-with shell`
See more here: https://www.vagrantup.com/docs/cli/up


### Plugins

- `vagrant plugin install vagrant-digitalocean` enable digitalocean
- `vagrant plugin install vagrant-scp`
- `vagrant plugin install vagrant-vbguest`
- `vagrant plugin install vagrant-reload`
- `vagrant plugin install vagrant-env` load .env files
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

## Synced folders
To sync folders in into the VM add following line
```ruby
config.vm.synced_folder ".", "/vagrant", type: "rsync" # "rsync" for digitalocean. "virtualbox" for virtualbox
```

Hidden files (dot files) are not synced. So we need to manually fetch them:
```ruby
config.vm.define "webserver", primary: true do |server|
      server.vm.provision "file", source: ".env", destination: ".env"
```

## Configuration
Create .env in root folder and copy the contents of .env.example. Fill out the fields with your information. The keys DIGITAL_OCEAN_TOKEN and SSH_KEY_NAME are configured from digital ocean. 


## Running api tests
in /python_api_tests: ```python3 minitwit_simulator.py "http://67.207.75.163:8080/api"```


## Reading log file on server
run command in vagrant ```tail -f -n100 out.log```