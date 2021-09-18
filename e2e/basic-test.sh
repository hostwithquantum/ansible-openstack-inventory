#!/usr/bin/env bash


source ./assert.sh

result_json="$(pwd)/result.json"

# This relies on ENV vars to be present. It's quick and dirty.

cd ../ && make e2e

assert_lt 0 $(cat ${result_json}|jq -r '._meta.hostvars|length')
if [ "$?" == 0 ]; then
  log_success "We found at least one host"
else
  log_failure "No host found"
fi

assert_true $(cat ${result_json}|jq -r '.docker_swarm_manager.hosts|length > 0')
if [ "$?" == 0 ]; then
  log_success "We found at least one host in docker_swarm_manager"
else
  log_failure "We found no host in docker_swarm_manager"
fi

assert_true $(cat ${result_json}|jq -r '.["portainer-agent-node"].children|length > 0')
if [ "$?" == 0 ]; then
  log_success "We found child groups in 'portainer-agent-node'"
else
  log_failure "Missing child groups in 'portainer-agent-node'"
fi

assert_true $(cat ${result_json}|jq -r '.["portainer-agent-node"].children[0] == "docker_swarm_manager"')
if [ "$?" == 0 ]; then
  log_success "'docker_swarm_manager' is child of portainer-agent-node"
else
  log_failure "'docker_swarm_manager' is not a child of portainer-agent-node"
fi

assert_true $(cat ${result_json}|jq -r '.["portainer-agent-node"].children[1] == "docker_swarm_worker"')
if [ "$?" == 0 ]; then
  log_success "'docker_swarm_worker' is child of portainer-agent-node"
else
  log_failure "'docker_swarm_worker' is not a child of portainer-agent-node"
fi

assert_true $(cat ${result_json}|jq -r '._meta.hostvars|to_entries|.[0].value.floating_ip|length > 0')
if [ "$?" == 0 ]; then
  log_success "The first host should have a FIP property"
else
  log_failure "The first host doesn't have a FIP property"
fi

assert_true $(cat ${result_json}|jq -r '._meta.hostvars|to_entries|.[].value.ansible_host|length > 0')
if [ "$?" == 0 ]; then
  log_success "Each host has a property ansible_host"
else
  log_failure "Missing property ansible_host (in hostvars)"
fi
