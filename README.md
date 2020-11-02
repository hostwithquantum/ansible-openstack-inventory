# ansible-openstack-inventory

Inventory script for a dynamic OpenStack-based inventory.

This greats a default group `all` and various child groups which are configured via `config.ini`. We also use multiple networks on instancs, therefor a default/access network for Ansible has to be cobnfigured.

## Usage

```
$ QUANTUM_CUSTOMER=... ./ansible-openstack-inventory
```

## Todo

- [x] return correct JSON format of hosts
- [x] how to add additional groups
- [x] implement `--list`
- [x] implement `--host node`
- [ ] better error handling (instead of `os.Exit`)

