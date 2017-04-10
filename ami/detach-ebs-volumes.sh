#!/bin/sh
set -e
set -o pipefail
# set -x

DEVICE=/dev/xvdcz

sudo yum install -y aws-cli jq

sudo umount -d "$DEVICE" || true

instance_id=$(curl -sS http://169.254.169.254/latest/meta-data/instance-id)
volume_id=$(
  aws ec2 describe-volumes --output json --filters "Name=attachment.instance-id,Values=${instance_id}" "Name=attachment.device,Values=${DEVICE}" \
  | jq -Mr '.Volumes[0].VolumeId'
)

aws ec2 detach-volume --volume-id "${volume_id}" --force
