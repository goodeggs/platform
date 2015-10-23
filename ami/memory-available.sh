#!/bin/bash
 
HOSTNAME="${COLLECTD_HOSTNAME:-`hostname -f`}"
INTERVAL="${COLLECTD_INTERVAL:-10}"
 
while sleep "$INTERVAL"; do

  available_mem_kb=$(cat /proc/meminfo | grep '^MemAvailable:' | awk -F ':' '{print $2}' | awk '{print $1}')
  total_mem_kb=$(cat /proc/meminfo | grep '^MemTotal:' | awk -F ':' '{print $2}' | awk '{print $1}')
  available_mem_b=$(( $available_mem_kb * 1024 ))
  available_mem_pct=$(printf "%.2f" $(awk "BEGIN { print ($available_mem_kb / $total_mem_kb) * 100 }"))

  echo "PUTVAL \"$HOSTNAME/memory/memory-available\" interval=$INTERVAL N:$available_mem_b"
  echo "PUTVAL \"$HOSTNAME/memory/percent-available\" interval=$INTERVAL N:$available_mem_pct"

done
