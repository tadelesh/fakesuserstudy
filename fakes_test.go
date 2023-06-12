//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fakesuserstudy

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/stretchr/testify/require"
)

// docs https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5@v5.1.0-beta.1#readme-fakes
// see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore@v1.7.0-beta.2/fake for general docs on fakes

func Test_VirtualMachinesClient_Get(t *testing.T) {
	// write a fake for VirtualMachinesClient.Get that includes the following data

	const (
		// the name of the VM returned from VirtualMachinesClient.Get
		vmName = "virtualmachine1"

		// the resource ID of the VM returned from VirtualMachinesClient.Get
		resourceID = "/fake/resource/id"
	)

	// TODO: write fake here

	// TODO: create client and connect it to the fake
	var client armcompute.VirtualMachinesClient

	vm, err := client.Get(context.Background(), "fake-resource-group", vmName, nil)

	// the result must satisfy the following conditions
	require.NoError(t, err)
	require.NotNil(t, vm.Name)
	require.EqualValues(t, vmName, *vm.Name)
	require.NotNil(t, vm.ID)
	require.EqualValues(t, resourceID, *vm.ID)
}

func Test_VirtualMachinesClient_Get_error(t *testing.T) {
	// write a fake for VirtualMachinesClient.Get that includes the following data

	const (
		// the HTTP status code of the failed request
		httpError = http.StatusBadRequest

		// the error code of the failed request
		errorCode = "ErrorResourceNotFound"
	)

	// TODO: write fake here

	// TODO: create client and connect it to the fake
	var client armcompute.VirtualMachinesClient

	vm, err := client.Get(context.Background(), "fake-resource-group", "virtualmachine1", nil)

	// the result must satisfy the following conditions
	require.Zero(t, vm)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, httpError, respErr.StatusCode)
	require.EqualValues(t, errorCode, respErr.ErrorCode)
}

func Test_VirtualMachinesClient_BeginCreateOrUpdate(t *testing.T) {
	// write a fake for VirtualMachinesClient.BeginCreateOrUpdate that includes the following data

	const (
		// the name of the VM returned when the long-running operation completes
		vmName = "virtualmachine1"

		// the resource ID of the VM returned when the long-running operation completes
		resourceID = "/fake/resource/id"
	)

	// TODO: write fake here. the poller must include two non-terminal responses

	// TODO: create client and connect it to the fake
	var client armcompute.VirtualMachinesClient

	poller, err := client.BeginCreateOrUpdate(context.Background(), "fake-resource-group", vmName, armcompute.VirtualMachine{}, nil)
	require.NoError(t, err)

	pollingIterations := 0

	for !poller.Done() {
		resp, err := poller.Poll(context.Background())
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, resp.StatusCode)

		pollingIterations++
	}

	require.EqualValues(t, 2, pollingIterations)

	vm, err := poller.Result(context.Background())

	// the result must satisfy the following conditions
	require.NoError(t, err)
	require.NotNil(t, vm.Name)
	require.EqualValues(t, vmName, *vm.Name)
	require.NotNil(t, vm.ID)
	require.EqualValues(t, resourceID, *vm.ID)
}

func Test_VirtualMachinesClient_NewListPager(t *testing.T) {
	// write a fake for VirtualMachinesClient.NewListPager that returns a total of
	// five VMs spread over two pages. the first page should include three VMs and
	// the second page should contain two VMs.
	// to keep things simple, the returned armcompute.VirtualMachine instances can be empty.

	// TODO: write fake here

	// TODO: create client and connect it to the fake
	var client armcompute.VirtualMachinesClient

	pager := client.NewListPager("fake-resource-group", nil)

	pageCount := 0
	vmCount := 0

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		pageCount++
		vmCount += len(page.Value)
	}

	// the results must satisfy the following conditions
	require.EqualValues(t, 2, pageCount)
	require.EqualValues(t, 5, vmCount)
}
