#!/bin/bash

# Make em quiet
pushd () {
  command pushd "$@" > /dev/null
}

popd () {
  command popd "$@" > /dev/null
}

# Exclude containers we don't want effected by limits
EXCLUDES_PATTERN=$(cat <<'EOF' | xargs | sed 's/ /|/g'
amazon/amazon-ecs-agent
goodeggs/logspout-goodeggs
yelp/docker-custodian
convox/api
convox/agent
goodeggs/ranch-api
EOF
)

# Build a list of stuff we DO want to effect
TARGETS=$( docker ps --no-trunc --format '{{.ID}} {{.Image}}' | grep -Ev "$EXCLUDES_PATTERN" | awk '{ print $1; }' | xargs)

# Apply Modifications
for a in $TARGETS; do

  # Disable swap (set memory and swap to same size)
  pushd /cgroup/memory/docker/${a}
  cat "./memory.limit_in_bytes" > "./memory.memsw.limit_in_bytes"
  popd

  # Keeping here for now incase we leave it later
  # Limit IOPS for all devices (5 IOPS/Sec to prevent thrashing, reads and writes)
  #pushd /cgroup/blkio/docker/$a
  #DEVICES="$(cat ./blkio.throttle.io_service_bytes  | awk '{print $1;}' | uniq | grep -v Total)"
  #for b in $DEVICES; do
  #  echo "${b} 5" > ./blkio.throttle.write_iops_device
  #  echo "${b} 5" > ./blkio.throttle.read_iops_device
  #done
  #popd
done

# No more swap
# /sbin/swapoff -a

# Log major pagefaults
for a in $TARGETS; do
  pushd /cgroup/memory/docker/$a
  RESULT=$(echo "id=${a} $(docker inspect --format 'image={{.Config.Image}} StartedAt="{{.State.StartedAt}}"' $a) pgmajfault=$(cat memory.stat | grep total_pgmajfault | awk '{print $2;}')")
  logger $RESULT
  echo $RESULT
  popd
done

