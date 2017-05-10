#!/bin/sh

for a in $(find /cgroup/memory/docker -type d | grep -v '/cgroup/memory/docker$'); do
  cat "$a/memory.limit_in_bytes" > "$a/memory.memsw.limit_in_bytes"
done
