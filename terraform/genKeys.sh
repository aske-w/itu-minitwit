#!/bin/bash


mkdir ./ssh_key #&& ssh-keygen -t rsa -b 4096 -q -N '' -f ./ssh_key/terraform
ssh-keygen -t rsa -m PEM -f ./ssh_key/terraform