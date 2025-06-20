package virtualmachine

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/kubevirt/terraform-provider-kubevirt/ci-tests/common"
	"github.com/kubevirt/terraform-provider-kubevirt/ci-tests/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

func TestVirtualMachineNetworking(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Machine Networking Suite")
}

var _ = Describe("Virtual Machine networking", func() {
	var (
		testDir   string
		testID    string
		namespace string
	)

	BeforeEach(func() {
		var err error
		if testDir, err = ioutil.TempDir("", "vm-network-test-"); err != nil {
			Fail(fmt.Sprintf("failed to create temp dir for terraform execution, with error: %s", err))
		}
		testID = uuid.New()
		namespace = fmt.Sprintf("vm-network-test-namespace-%s", testID)
		common.CreateNamespace(namespace)
	})

	AfterEach(func() {
		common.DeleteNamespace(namespace)
		os.RemoveAll(testDir)
	})

	Context("With pod network", func() {
		It("should create a VM with pod network", func() {
			testName := "vm-pod-network"
			tfExecPath := "terraform"
			
			// Set up Terraform
			err := terraform.Init(testDir, testName, tfExecPath)
			Expect(err).NotTo(HaveOccurred(), "Failed to init terraform")
			
			// Create VM with pod network
			podNetworkVM := fmt.Sprintf(`
provider "kubevirt" {}

resource "kubevirt_virtual_machine" "test-vm" {
  metadata {
    name      = "test-vm-pod-network"
    namespace = "%s"
    labels = {
      "key1" = "value1"
    }
  }
  spec {
    run_strategy = "Always"
    template {
      metadata {
        labels = {
          "kubevirt.io/vm" = "test-vm-pod-network"
        }
      }
      spec {
        domain {
          resources {
            requests = {
              memory = "64M"
              cpu    = 1
            }
          }
          devices {
            interface {
              name                     = "default"
              interface_binding_method = "InterfaceBridge"
            }
          }
        }
        network {
          name = "default"
          network_source {
            pod {
              vm_network_cidr = "10.0.0.0/24"
            }
          }
        }
      }
    }
  }
}`, namespace)

			// Write the config to a file
			err = ioutil.WriteFile(fmt.Sprintf("%s/main.tf", testDir), []byte(podNetworkVM), 0644)
			Expect(err).NotTo(HaveOccurred(), "Failed to write Terraform config")

			// Apply the config
			var applyOpts []tfexec.ApplyOption
			err = terraform.Apply(testDir, tfExecPath, nil, applyOpts...)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply Terraform config")

			// Verify VM exists
			common.ValidateVirtualMachine("test-vm-pod-network", namespace, nil)

			// Destroy the infrastructure
			var destroyOpts []tfexec.DestroyOption
			err = terraform.Destroy(testDir, tfExecPath, destroyOpts...)
			Expect(err).NotTo(HaveOccurred(), "Failed to destroy Terraform infrastructure")

			// Verify VM is deleted
			common.ValidateVirtualMachine("test-vm-pod-network", namespace, nil)
		})
	})

	Context("With multus network", func() {
		It("should create a VM with multus network", func() {
			testName := "vm-multus-network"
			tfExecPath := "terraform"
			
			// Set up Terraform
			err := terraform.Init(testDir, testName, tfExecPath)
			Expect(err).NotTo(HaveOccurred(), "Failed to init terraform")
			
			// Create VM with multus network
			multusNetworkVM := fmt.Sprintf(`
provider "kubevirt" {}

resource "kubevirt_virtual_machine" "test-vm" {
  metadata {
    name      = "test-vm-multus-network"
    namespace = "%s"
    labels = {
      "key1" = "value1"
    }
  }
  spec {
    run_strategy = "Always"
    template {
      metadata {
        labels = {
          "kubevirt.io/vm" = "test-vm-multus-network"
        }
      }
      spec {
        domain {
          resources {
            requests = {
              memory = "64M"
              cpu    = 1
            }
          }
          devices {
            interface {
              name                     = "multus-network"
              interface_binding_method = "InterfaceBridge"
            }
          }
        }
        network {
          name = "multus-network"
          network_source {
            multus {
              network_name = "default/multus-test"
              default      = true
            }
          }
        }
      }
    }
  }
}`, namespace)

			// Write the config to a file
			err = ioutil.WriteFile(fmt.Sprintf("%s/main.tf", testDir), []byte(multusNetworkVM), 0644)
			Expect(err).NotTo(HaveOccurred(), "Failed to write Terraform config")

			// Apply the config
			var applyOpts []tfexec.ApplyOption
			err = terraform.Apply(testDir, tfExecPath, nil, applyOpts...)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply Terraform config")

			// Verify VM exists
			common.ValidateVirtualMachine("test-vm-multus-network", namespace, nil)

			// Destroy the infrastructure
			var destroyOpts []tfexec.DestroyOption
			err = terraform.Destroy(testDir, tfExecPath, destroyOpts...)
			Expect(err).NotTo(HaveOccurred(), "Failed to destroy Terraform infrastructure")

			// Verify VM is deleted
			common.ValidateVirtualMachine("test-vm-multus-network", namespace, nil)
		})
	})

	Context("With both pod and multus networks", func() {
		It("should create a VM with both pod and multus networks", func() {
			testName := "vm-dual-network"
			tfExecPath := "terraform"
			
			// Set up Terraform
			err := terraform.Init(testDir, testName, tfExecPath)
			Expect(err).NotTo(HaveOccurred(), "Failed to init terraform")
			
			// Create VM with both pod and multus networks
			dualNetworkVM := fmt.Sprintf(`
provider "kubevirt" {}

resource "kubevirt_virtual_machine" "test-vm" {
  metadata {
    name      = "test-vm-dual-network"
    namespace = "%s"
    labels = {
      "key1" = "value1"
    }
  }
  spec {
    run_strategy = "Always"
    template {
      metadata {
        labels = {
          "kubevirt.io/vm" = "test-vm-dual-network"
        }
      }
      spec {
        domain {
          resources {
            requests = {
              memory = "64M"
              cpu    = 1
            }
          }
          devices {
            interface {
              name                     = "default"
              interface_binding_method = "InterfaceBridge"
            }
            interface {
              name                     = "multus-network"
              interface_binding_method = "InterfaceBridge"
            }
          }
        }
        network {
          name = "default"
          network_source {
            pod {
              vm_network_cidr = "10.0.0.0/24"
            }
          }
        }
        network {
          name = "multus-network"
          network_source {
            multus {
              network_name = "default/multus-test"
              default      = true
            }
          }
        }
      }
    }
  }
}`, namespace)

			// Write the config to a file
			err = ioutil.WriteFile(fmt.Sprintf("%s/main.tf", testDir), []byte(dualNetworkVM), 0644)
			Expect(err).NotTo(HaveOccurred(), "Failed to write Terraform config")

			// Apply the config
			var applyOpts []tfexec.ApplyOption
			err = terraform.Apply(testDir, tfExecPath, nil, applyOpts...)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply Terraform config")

			// Verify VM exists
			common.ValidateVirtualMachine("test-vm-dual-network", namespace, nil)

			// Destroy the infrastructure
			var destroyOpts []tfexec.DestroyOption
			err = terraform.Destroy(testDir, tfExecPath, destroyOpts...)
			Expect(err).NotTo(HaveOccurred(), "Failed to destroy Terraform infrastructure")

			// Verify VM is deleted
			common.ValidateVirtualMachine("test-vm-dual-network", namespace, nil)
		})
	})
})
