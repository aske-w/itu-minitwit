Vagrant.configure("2") do |config|
    config.vm.box = "generic/ubuntu1804"
    #config.ssh.private_key_path = '~/.ssh/id_rsa'

    #config.vm.network "private_network", type: "dhcp"

    # For two way synchronization you might want to try `type: "virtualbox"`
    config.vm.synced_folder ".", "/vagrant", type: "virtualbox" # was rsync before
  
    # config.vm.define "dbserver", primary: true do |server|
    #   server.vm.network "private_network", ip: "192.168.56.2"
    #   # config.vm.network "forwarded_port", guest: 27017, host: 37017
    #   # config.vm.network "forwarded_port", guest: 28017, host: 38017
    #   server.vm.provider "virtualbox" do |vb|
    #     vb.memory = "1024"
    #   end
    #   server.vm.hostname = "dbserver"
    #   server.vm.provision "shell", privileged: false, inline: <<-SHELL
    #       echo "Installing MongoDB"
    #       wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | sudo apt-key add -
    #       echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.2.list
    #       sudo apt-get update
    #       sudo apt-get install -y mongodb-org
    #       sudo mkdir -p /data/db
    #       sudo sed -i '/  bindIp:/ s/127.0.0.1/0.0.0.0/' /etc/mongod.conf
    #       sudo systemctl start mongod
    #       mongorestore --gzip /vagrant/dump
    #   SHELL

    config.vm.define "webserver", primary: true do |server|
        #server.vm.network "private_network", ip: "192.168.68.35"
        
        #config.vm.provision "file", source: "~/path/to/host/folder", destination: "$HOME/remote/newfolder"
        #server.vm.network "forwarded_port", guest: 5000, host: 7831
        server.vm.provider "virtualbox" do |vb|
          vb.memory = "1024"
        end
        server.vm.hostname = "webserver"
        server.vm.provision "shell", privileged: false, inline: <<-SHELL
            echo "INSIDE PROVISION SCRIPT!"
            export GOPATH=$HOME/go
            export GO_VERSION="go1.17.7.linux-amd64"
            sudo curl -O https://storage.googleapis.com/golang/$GO_VERSION.tar.gz
            sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_VERSION.tar.gz

            


            echo ". $HOME/.bashrc" >> $HOME/.bash_profile
            echo "export PATH=/usr/local/go/bin:$PATH" >> $HOME/.bash_profile
            export PATH="/usr/local/go/bin:$PATH"
            source $HOME/.bash_profile

            cp -r /vagrant/* $HOME
            nohup go run main.go
            
            # export DB_IP="192.168.56.2"
            # echo "Installing Anaconda..."
            # sudo wget https://repo.anaconda.com/archive/Anaconda3-2019.07-Linux-x86_64.sh -O $HOME/Anaconda3-2019.07-Linux-x86_64.sh
        
            # bash $HOME/Anaconda3-2019.07-Linux-x86_64.sh -b
            
            # echo ". $HOME/.bashrc" >> $HOME/.bash_profile
            # echo "export PATH=$HOME/anaconda3/bin:$PATH" >> $HOME/.bash_profile
            # export PATH="$HOME/anaconda3/bin:$PATH"
            # rm $HOME/Anaconda3-2019.07-Linux-x86_64.sh
            # source $HOME/.bash_profile
            # pip install Flask-PyMongo
            # cp -r /vagrant/* $HOME
            # nohup python minitwit.py > out.log 2>&1 &
            # echo "================================================================="
            # echo "=                            DONE                               ="
            # echo "================================================================="
            # echo "Navigate in your browser to:"
            # echo "http://192.168.56.3:5000"
        SHELL
      end
    
      #config.vm.provision "shell", privileged: false, inline: <<-SHELL
      #  sudo apt-get update
      #SHELL
    end