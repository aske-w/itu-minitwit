cat test_env.sh ufw_script.sh worker_script.sh > worker_provision.sh
cat test_env.sh ufw_script.sh manager_script.sh > manager_provision.sh

chmod +x worker_provision.sh
chmod +x manager_provision.sh