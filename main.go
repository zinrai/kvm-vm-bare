package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	imageDir = "/var/lib/libvirt/images"
)

func main() {
	// Check if required commands exist
	checkRequiredCommands()

	// Define flags
	diskSize := flag.String("size", "20G", "Size of the virtual disk")
	memory := flag.Int("memory", 1024, "Memory in MB")
	vcpus := flag.Int("vcpus", 1, "Number of virtual CPUs")
	network := flag.String("network", "network=default", "Network configuration (e.g., 'bridge=br0' or 'network=default')")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] VM_NAME\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check if VM name is provided
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Error: VM name is required as the last argument")
		flag.Usage()
		os.Exit(1)
	}
	vmName := args[0]

	// Create disk image
	diskPath := filepath.Join(imageDir, fmt.Sprintf("%s.qcow2", vmName))
	createDisk(diskPath, *diskSize)

	// Generate XML
	xml := generateXML(vmName, diskPath, *memory, *vcpus, *network)

	// Define VM
	defineVM(xml)

	fmt.Printf("Empty VM '%s' created successfully with network: %s\n", vmName, *network)
	fmt.Printf("Disk image created at: %s\n", diskPath)
}

func checkRequiredCommands() {
	requiredCommands := []string{"qemu-img", "virsh", "virt-install"}

	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			log.Fatalf("Required command not found: %s\nPlease install the necessary packages (usually libvirt-clients, libvirt-daemon-system, and qemu-utils)", cmd)
		}
	}
}

func createDisk(path, size string) {
	// Check if the directory exists
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", filepath.Dir(path))
	}

	cmd := exec.Command("sudo", "qemu-img", "create", "-f", "qcow2", path, size)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to create disk: %v\nOutput: %s", err, output)
	}
}

func generateXML(name, diskPath string, memory, vcpus int, network string) string {
	args := []string{
		"virt-install",
		"--name", name,
		"--memory", fmt.Sprintf("%d", memory),
		"--vcpus", fmt.Sprintf("%d", vcpus),
		"--disk", fmt.Sprintf("path=%s,format=qcow2", diskPath),
		"--os-variant", "generic",
		"--print-xml",
	}

	// Validate and add network configuration
	if strings.HasPrefix(network, "bridge=") || strings.HasPrefix(network, "network=") {
		args = append(args, "--network", network)
	} else {
		log.Fatalf("Invalid network configuration: %s. Use 'bridge=BRIDGE' or 'network=NAME'", network)
	}

	cmd := exec.Command("sudo", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to generate XML: %v", err)
	}

	return out.String()
}

func defineVM(xml string) {
	cmd := exec.Command("sudo", "virsh", "define", "/dev/stdin")
	cmd.Stdin = strings.NewReader(xml)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to define VM: %v\nOutput: %s", err, output)
	}
}
