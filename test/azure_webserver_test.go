package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "b40c6add-0fa1-4055-9718-f3d28641e566"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "marc",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	imagePublisher := terraform.Output(t, terraformOptions, "vm_image_publisher")
	imageOffer := terraform.Output(t, terraformOptions, "vm_image_offer")
	imageSKU := terraform.Output(t, terraformOptions, "vm_image_sku")
	imageVersion := terraform.Output(t, terraformOptions, "vm_image_version")
	// Expected values from Terraform configuration
	expectedPublisher := "Canonical"
	expectedOffer := "0001-com-ubuntu-server-jammy"
	expectedSKU := "22_04-lts-gen2"
	expectedVersion := "latest"
	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is connected to the VM
	nicExists := azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID)
	t.Logf("Checking if NIC %s exists in resource group %s", nicName, resourceGroupName)
	assert.True(t, nicExists, "Network Interface should exist")

	// Validate the expected image details
	assert.Equal(t, expectedPublisher, imagePublisher, "VM should be from Canonical")
	assert.Equal(t, expectedOffer, imageOffer, "VM offer should be 0001-com-ubuntu-server-jammy")
	assert.Equal(t, expectedSKU, imageSKU, "VM should be running Ubuntu 22.04 LTS")
	assert.Equal(t, expectedVersion, imageVersion, "VM verison should be latest")
}
	
