# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
	config.vm.box = "digital_ocean"
	config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
	config.ssh.private_key_path = './minitwit_dd'
	#config.vm.network "private_network", type: "dhcp"

	# For two way synchronization you might want to try `type: "virtualbox"`
	config.vm.synced_folder "remote_files", "/vagrant", type: "rsync" # was rsync before
	config.env.enable

	config.vm.define "webserver", primary: true do |server|
		# server.vm.network "private_network", ip: "192.168.62.1"
		# server.vm.network "forwarded_port", guest: 8080, host: 8080
	
		#config.vm.provision "file", source: "~/path/to/host/folder", destination: "$HOME/remote/newfolder"
		server.vm.provider "digital_ocean" do |provider|
			provider.ssh_key_name = 'do_ssh_key'
			provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
			provider.image = 'docker-18-04'
			provider.region = 'fra1'
			provider.size = 's-1vcpu-1gb'
			provider.privatenetworking = true
		end
		
		server.vm.hostname = "webserver"
		server.vm.provision "shell", inline: <<-SHELL
			echo "INSIDE PROVISION SCRIPT!"

			echo -e "\nOpening port for minitwit ...\n"
			ufw allow 5000

			echo -e "\nOpening port for minitwit ...\n"
			echo ". $HOME/.bashrc" >> $HOME/.bash_profile

			# export THIS_IP=`hostname -I | cut -d" " -f1`
			
			# # stops the current running process
			# if [[ -e save_pid.txt ]]; then
			#   PID=`cat save_pid.txt`
				
			#   # If the process is running - kill it
			#   echo "Killing old running process $PID"
			#   kill -0 $PID

			#   rm save_pid.txt  
			# fi
				
			# # build executable
			# go build main.go
			# # run in background, while logging to out.log
			# nohup ./main > out.log 2>&1 & echo $! > save_pid.txt 
			# # "dollarsign exclamation" is PID of last program (./main)
			# echo "http://${THIS_IP}:8080"
				
		SHELL
	end
	
	# config.vm.provision "shell", privileged: false, inline: <<-SHELL
	#  sudo apt-get update
	# SHELL
  end