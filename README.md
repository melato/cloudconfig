cloudconfig - provides Go structures and interfaces for reading and applying
[cloud-init](https://cloud-init.io/) configuration files.

# Supported cloud-init features
The following cloud-init modules (sections) are supported and applied in this order:
- packages
- write_files
- users
- runcmd
- write_files with defer: true

The files are in YAML format and their first line must be "#cloud-config".

# Examples
- [examples/](https://github.com/melato/cloudconfig/blob/main/examples/)
- [Examples in the cloud-init documentation](https://cloudinit.readthedocs.io/en/latest/reference/examples.html)

# sudo vs doas
The users sudo option configures either sudo or Alpine doas,
according to which directory it finds, /etc/opt/sudo.d, or /etc/opt/doas.d.

doas and sudo configuration strings are not compatible.
If you specify "sudo: true", then an appropriate configuration can be applied to both sudo and doas.
But if you specify a string, then it should be either sudo or doas specific.
It will not work correctly if both sudo.d and doas.d exist.

# Implementations
This project provides a local implementation that applies the cloud-config files to the local machine.

Other projects provide implementations for applying cloud-config files to other systems:
- [cloudconfiglxd](https://github.com/melato/cloudconfiglxd)
applies cloud-config files to LXD instances, via the LXD InstanceServer API.

- [cloudconfigincus](https://github.com/melato/cloudconfigincus)
applies cloud-config files to Incus instances, via the Incus InstanceServer API.


# Standalone executable
main/cloudconfig.go can be used to compile a standalone executable with the local implementation.
It may also be useful for examining and debugging cloud-config files.

## Usage
```
cloudconfig apply [-os <ostype>] <cloud-config-file>...
```
ostype is needed to for packages and users, since different distributions have
different package systems and may have differences in how they create users.

Supported ostypes are: alpine, debian.  Others can be added.
  
## compile

```
cd main
git log -1 --format=%cd > version
# or: date > version
go install cloudconfig.go
```

# Limitations
- write_files supports only plain text, without any encoding.
- users supports name, uid, shell, homedir, no_create_home, primary_group, groups, gecos, ssh_authorized_keys, sudo.
- Applying a user to an existing user may not work as specified in cloud-init.

Different implementations may behave differently.

The goal is to adhere to a useful small subset of the cloud-init specification, which can be implemented easily.