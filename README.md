# kvm-vm-bare

`kvm-vm-bare` is a command-line tool that simplifies the process of creating empty virtual machines using KVM (Kernel-based Virtual Machine). This tool is designed for users who need to quickly set up base VMs without an operating system installed.

## Features

- Create empty VMs with customizable resources
- Flexible network configuration

## Prerequisites

Before using this tool, ensure you have the following installed on your system:

- sudo
- virsh
- qemu-img
- virt-install

## Installation

Build the tool:

```
$ go build
```

## Usage

The basic syntax for using the tool is:

```
$ ./kvm-vm-bare [OPTIONS] VM_NAME
```

### Examples:

Create a VM with default settings:

```
$ ./kvm-vm-bare myvm
```

Create a VM with custom resources:

```
$ ./kvm-vm-bare -size 30G -memory 2048 -vcpus 2 myvm
```

Create a VM with a specific bridge network:

```
$ ./kvm-vm-bare -network bridge=br0 myvm
```

## Network Configuration

You can specify the network configuration using the `-network` option. The value should be in the format accepted by virt-install's --network option. For example:

- `bridge=BRIDGE`: Connect to a bridge device
- `network=NAME`: Connect to a virtual network

## Notes

- This tool requires sudo privileges to run as it needs to create disk images and define VMs.
- The created VMs are empty and do not have an operating system installed. You'll need to manually install an OS after creation.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
