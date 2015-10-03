#!/bin/sh

# NOTE: this exists because the ECS AMI does not include scp.
# we can kill this and the hacky sleep 30 off once https://github.com/mitchellh/packer/pull/2504 is merged.

yum install -y openssh-clients
