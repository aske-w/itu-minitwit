Vagrant.configure("2") do |config|
  config.vm.box = "digital_ocean"
  config.ssh.private_key_path = '~/.ssh/id_rsa'
  config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
#   config.vm.network "private_network", type: "dhcp"

  # For two way synchronization you might want to try `type: "virtualbox"`
  config.vm.synced_folder ".", "/vagrant", type: "rsync" # was rsync before
  config.env.enable

  config.vm.define "webserver", primary: true do |server|
    
      # server.vm.network "private_network", ip: "192.168.62.1"
      # server.vm.network "forwarded_port", guest: 8080, host: 8080
      
      #config.vm.provision "file", source: "~/path/to/host/folder", destination: "$HOME/remote/newfolder"
      server.vm.provider "virtualbox" do |vb|
        provider.ssh_key_name = ENV["SSH_KEY_NAME"]
        provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
        provider.image = 'ubuntu-18-04-x64'
        provider.region = 'fra1'
        provider.size = 's-1vcpu-1gb'
        provider.privatenetworking = true
      end

      
      server.vm.hostname = "webserver"
      server.vm.provision "file", source: ".env", destination: ".env"
      server.vm.provision "shell", privileged: false, inline: <<-SHELL
          echo "INSIDE PROVISION SCRIPT!"
          
          export GO_VERSION="go1.17.7.linux-amd64"
          sudo curl -O https://storage.googleapis.com/golang/$GO_VERSION.tar.gz
          sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_VERSION.tar.gz

          

          echo ". $HOME/.bashrc" >> $HOME/.bash_profile
          echo "export PATH=/usr/local/go/bin:$PATH" >> $HOME/.bash_profile
          export PATH="/usr/local/go/bin:$PATH"
          source $HOME/.bash_profile

          cp -r /vagrant/* $HOME
          export THIS_IP=`hostname -I | cut -d" " -f1`
          echo "http://${THIS_IP}:8080"
          nohup go run main.go
          
         
      SHELL
    end
  
    config.vm.provision "shell", privileged: false, inline: <<-SHELL
     sudo apt-get update
    SHELL
  end