//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fakesuserstudy

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/stretchr/testify/require"
)

func Test_VirtualMachinesClient_Get(t *testing.T) {
	// write a fake for VirtualMachinesClient.Get that satisfies the following requirements

	const (
		vmName            = "virtualmachine1"
		resourceGroupName = "fake-resource-group"
	)

	// the fake VM must return the provided name and its ID contain the provided resource group name.

	// TODO: populate vm with response from fake
	var vm armcompute.VirtualMachine

	// the returned VM must satisfy the following conditions
	require.NotNil(t, vm.Name)
	require.Equal(t, vmName, *vm.Name)
	require.NotNil(t, vm.ID)
	require.Contains(t, *vm.ID, resourceGroupName)
}

func Test_VirtualMachinesClient_BeginDelete(t *testing.T) {
	// write a fake for VirtualMachinesClient.BeginDelete that satisfies the following requirements

	const (
		vmName            = "virtualmachine1"
		resourceGroupName = "fake-resource-group"
	)

	// TODO: populate pollingErr with the error after polling completes.
	// the fake should include at least one non-terminal response.
	var pollingErr error

	// the LRO must terminate in a way to satisfy the following conditions
	require.Error(t, pollingErr)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, pollingErr, &respErr)
}

func Test_VirtualMachinesClient_NewListPager(t *testing.T) {
	// write a fake for VirtualMachinesClient.NewListPager that satisfies the following requirements

	const (
		resourceGroupName = "fake-resource-group"
	)

	// the fake must return a total of five VMs over two pages.
	// to keep things simple, the returned armcompute.VirtualMachine instances can be empty.
	// while iterating over pages, the fake must return one transient error before the final page

	// TODO: populate vmCount with the number of VMs returned
	var vmCount int

	// TODO: populate pageCount with the number of returned pages
	var pageCount int

	// TODO: populate errCount with the number of transient errors encountered during paging
	var errCount int

	// the results must satisfy the following conditions
	require.Equal(t, 5, vmCount)
	require.Equal(t, 2, pageCount)
	require.Equal(t, 1, errCount)
}
