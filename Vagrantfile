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
      server.vm.provider "digital_ocean" do |provider|
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
          
         

          sudo apt-get update
          
          # SQLITE is not installed:
          if [[ ! -e `which sqlite3` ]]; then
            sudo apt-get install sqlite3
          else 
            echo "Sqlite3 already installed. Skipping."
          fi


          # GCC is not installed:
          if [[ ! -e `which gcc` ]]; then
            sudo apt-get install -y build-essential
          else 
            echo "Gcc already installed. Skipping."
          fi
          # Go is not installed:
          if [[ ! -e `which go` ]]; then
            export GO_VERSION="go1.17.7.linux-amd64"
            sudo curl -O https://storage.googleapis.com/golang/$GO_VERSION.tar.gz
            sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_VERSION.tar.gz
            
            
            
            echo ". $HOME/.bashrc" >> $HOME/.bash_profile
            echo "export PATH=/usr/local/go/bin:$PATH" >> $HOME/.bash_profile
            export PATH="/usr/local/go/bin:$PATH"
            source $HOME/.bash_profile
          else 
            echo "Go already installed. Skipping."
          fi  

          # dont overwrite the db in the VM
          if [[ -e $HOME/db.db ]]; then
            rm /vagrant/db.db
          else
            # cat schema.sql | sqlite3 db.db
            echo "No DB file found. Overwriting."
          fi 

          cp -r /vagrant/* $HOME
          export THIS_IP=`hostname -I | cut -d" " -f1`
          
          # stops the current running process
          if [[ -e save_pid.txt ]]; then
            PID=`cat save_pid.txt`
            
            # If the process is running - kill it
            echo "Killing old running process $PID"
            kill -0 $PID

            rm save_pid.txt  
          fi
              
          # build executable
          go build main.go
          # run in background, while logging to out.log
          nohup ./main > out.log 2>&1 & echo $! > save_pid.txt 
          # "dollarsign exclamation" is PID of last program (./main)
          echo "http://${THIS_IP}:8080"
            
      SHELL
    end
  
    config.vm.provision "shell", privileged: false, inline: <<-SHELL
     sudo apt-get update
    SHELL
  end