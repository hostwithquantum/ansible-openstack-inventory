# ansible-openstack-inventory

Inventory script for a dynamic OpenStack-based inventory.

This script creates a default group `all` and various child groups which are configured via `config.ini`. We also use multiple networks on instances, therefor a default/access network for Ansible has to be configured.

## Project Status

This is not a general purpose inventory for all OpenStack clouds, but instead it's heavily tailored towards what we use and need to bootstrap nodes for our [application hosting service](https://www.planetary-quantum.com).

Specifically in regards to Docker Swarm, inventory script returns nodes and adds groups and labels for each node/group:

 - `docker_swarm_manager`
 - `docker_swarm_worker`

Group membership and (`swarm_`)labels are determined by instance metadata:

 - `com.planetary-quantum.meta.role` (`manager` (default), worker)
 - `com.planetary-quantum.meta.label`

Because documentation on dynamic inventories is a bit sparse, we decided to release this code to the broader community. And despite Ansible using Python, we wrote this inventory script in Go(lang) as we felt that at this part in the stack, we should run something less brittle.

So feel free to use, copy, fork and ad[a,o]pt - and feel free to contribute. Please keep in mind, that the code in this repository has to work for our use-case first.

## Usage

We use gophercloud to interface with OpenStack, the usual environment variables are listed in [`.envrc-dist`](.envrc-dist).

```
$ QUANTUM_CUSTOMER=... ./ansible-openstack-inventory
```

## Todo

- [x] return correct JSON format of hosts
- [x] how to add additional groups
- [x] implement `--list`
- [x] implement `--host node`
- [ ] better error handling (instead of `os.Exit`)

## License

BSD-2-Clause, Copyright 2021 (and beyond) [Planetary Quantum GmbH](https://www.planetary-quantum.com/service/impressum/)

