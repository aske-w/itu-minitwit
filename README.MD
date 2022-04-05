# Maintainability and Technical Debt estimation badges

## Sonarqube

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=aske-w_itu-minitwit&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=aske-w_itu-minitwit)


## Code Climate

[![Maintainability](https://api.codeclimate.com/v1/badges/305085e64353cb758288/maintainability)](https://codeclimate.com/github/aske-w/itu-minitwit/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/305085e64353cb758288/test_coverage)](https://codeclimate.com/github/aske-w/itu-minitwit/test_coverage)

## Better Code Hub

[![BCH compliance](https://bettercodehub.com/edge/badge/aske-w/itu-minitwit?branch=development)](https://bettercodehub.com/)

# Digital Ocean.

Floating IP (static ip): `67.207.75.163`

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

Make sure the `requests`library is installed. You can install it using `pip install requests`.
Also make sure the library is added to your path afterwards.
in /python_api_tests: `python3 minitwit_simulator.py "http://67.207.75.163:8080/api"`

## Reading log file on server

run command in vagrant `tail -f -n100 out.log`

## Creating tables in an empty db.db

Run the command `cat schema.sql | sqlite3 db.db` in bash

## Encrypting keys in travis

run cmd `travis encrypt-file minitwit_dd --com -r aske-w/itu-minitwit`

## Accessing the DB via MySQL CLI inside container

Run command `mysql -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE`

## Monitoring

docker-compose -d will run prometheus and grafana. If you are running the go server locally, configure the filed `targets` in the file `prometheues/prometheues.yml` to

```yaml
- targets:
    - host.docker.internal:8080
```

else specify localhost://8080.
In order to add data source in Grafana, in Grafana navigate to configuration/Data sources and press "Add Data Source". Select prometheues and set the `url` to the host of promeheues (http://localhost:9090). Then it's possible to create Dashboard with the different metrics.

### Creating Dashboard in Grafana

Once promotheues has been added as a data source we can under dashboard/new `add a new panel`. We can for example monitor the different container memory usage in bytes by this query `sum(container_memory_usage_bytes{name=~".+"}) by(name)`. It only shows containers with a name.

## Static analysis tools

On each pr or push to the `development` branch, the github action `static_analysis` is run. The action checks code formatting, constructs and potientiel security issues in source code. If a check fails the pipleline will fail.

### Format

Checks the code formatting. In case the code is not formatted correctly, it will throw an error and stop the pipeline.

### Examine source

Reports suspicious constructs in the source code. For example if `Printf` function arguments align with format string.

### Gosec Security Scanner

Scan source code for security problems and reports them. Gosec can be configured to run a set of rules. In our case we run all rules except for G104 `./bin/gosec -exclude=G104 ./...`. G104 audits errors not checked and since we do not handle all errors, the check will fail. However this should be handled in the future.

### Container scanning

After building the docker image https://github.com/Azure/container-scan this scans it for vulnerabilites. (cicd.yml)

## Elasticsearch

If you have issues with the container crashing, increase the container's max memory.

### .htpasswd

Run the script in `.htpasswd.example` and replace username with the username used for kibana and PASSWORD with a password

### Kibana

Navigate to port 5601. Under the logs tab type following command to query minitwit-server logs:
`container.image.name: "itu-minitwit_server"`

### Setup Kibana index pattern

On port 5601 navigate to `Stack management` at the bottom of the sidebar. Then check if the index `itu-minitwit-*` index has been registered under `data/Index Management`.
Under `Kibana/Index Patterns` do the following

-`Create new index pattern`.

- Type `itu-minitwit-*` as the name.

- Select the `@timestamp`

Once the index pattern is created navigate `Analytics/Discover` in the sidebar. The select the newly created index pattern.

We can the filter by fields. For example if we only want to see the message of the logs we can under the `Avaliable fields` press the plus icon to filter by messages.


## Docker swarm 

### manager
username: manager
password: 

### worker1
username: worker1
password: 

### worker2
username: worker2
password: 