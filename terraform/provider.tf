# api token
# here it is exported in the environment like
# export TF_VAR_do_token=xxx
variable "DOTOKEN" {}

# make sure to generate a pair of ssh keys
variable "pub_key" {}
variable "pvt_key" {}
variable "manager_tags" {}
variable "manager_coun" {}
variable "manager_size" {}
variable "manager_name" {}
variable "worker_tags" {}
variable "worker_coun" {}
variable "worker_size" {}
variable "worker_name" {}
variable "image" {}
variable "region" {}

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
  token = var.DOTOKEN
}