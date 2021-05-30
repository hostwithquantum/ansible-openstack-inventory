#!/usr/bin/env bash
result_json="$(pwd)/e2e-result.json"

# This relies on ENV vars to be present. It's quick and dirty.

make build && \
./dist/ansible-openstack-inventory_darwin_amd64/ansible-openstack-inventory --list > ${result_json} && \
# && cat ${result_json}|jq '.' \
echo "Must have at least one host ('> 0'):" && \
cat ${result_json}|jq -r '._meta.hostvars|length'

echo "Must have host in group 'docker_swarm_manager':" &&
cat ${result_json}|jq -r '.docker_swarm_manager.hosts|length > 0' &&

echo "Must have children in 'portainer-agent-node':" &&
cat ${result_json}|jq -r '.["portainer-agent-node"].children|length > 0' &&

echo "'docker_swarm_manager' is child of portainer-agent-node:" &&
cat ${result_json}|jq -r '.["portainer-agent-node"].children[0] == "docker_swarm_manager"' &&

echo "'docker_swarm_worker' is child of portainer-agent-node:" &&
cat ${result_json}|jq -r '.["portainer-agent-node"].children[1] == "docker_swarm_worker"' &&

# this is a stretch, but usually the case
echo "First host should be a manager and have a floating IP:" &&
cat ${result_json}|jq -r '._meta.hostvars|to_entries|.[0].value.floating_ip|length > 0' &&

echo "Each host should have 'ansible_host' variable": &&
cat ${result_json}|jq -r '._meta.hostvars|to_entries|.[].value.ansible_host|length > 0'
