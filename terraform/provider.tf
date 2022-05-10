# api token
# here it is exported in the environment like
# export TF_VAR_do_token=xxx
variable "do_token" {}

# make sure to generate a pair of ssh keys
variable "pub_key" {}
variable "pvt_key" {}
variable "manager_tags" {}
variable "manager_count" {}
variable "manager_size" {}
variable "manager_name" {}
variable "worker_tags" {}
variable "worker_count" {}
variable "worker_size" {}
variable "worker_name" {}
variable "image" {}
variable "region" {}

variable "state_file" {}
variable "manager_ip" {}
variable "bucket_name" {}
variable "access_key" {}
variable "secret_key" {}
variable "stack_name" {}
variable "prometheus_backup_interval" {}
variable "space_endpoint" {}

# setup the provider
terraform {
        required_providers {
                digitalocean = {
                        source = "digitalocean/digitalocean"
                        version = "~> 2.8.0"
                }
                null = {
                        source = "hashicorp/null"
                        version = "3.1.0"
                }
        }
}



provider "digitalocean" {
  token = var.do_token
}