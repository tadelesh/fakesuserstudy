//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fakesuserstudy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4/fake"
	"github.com/stretchr/testify/require"
)

type fakeVirtualMachinesServer struct {
	fake.VirtualMachinesServer
}

func (v *fakeVirtualMachinesServer) Get(ctx context.Context, resourceGroupName string, vmName string, options *armcompute.VirtualMachinesClientGetOptions) (resp fake.Responder[armcompute.VirtualMachinesClientGetResponse], err fake.ErrorResponder) {
	resp = fake.Responder[armcompute.VirtualMachinesClientGetResponse]{}
	resp.Set(armcompute.VirtualMachinesClientGetResponse{
		VirtualMachine: armcompute.VirtualMachine{
			Name: to.Ptr(vmName),
			ID:   to.Ptr(fmt.Sprintf("/subscriptions/subscriptionID/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s", resourceGroupName, vmName)),
		},
	})
	return
}

func Test_VirtualMachinesClient_Get(t *testing.T) {
	// write a fake for VirtualMachinesClient.Get that satisfies the following requirements

	const (
		vmName            = "virtualmachine1"
		resourceGroupName = "fake-resource-group"
	)

	// the fake VM must return the provided name and its ID contain the provided resource group name.

	// TODO: populate vm with response from fake
	client, err := armcompute.NewVirtualMachinesClient("subscriptionID", azfake.NewTokenCredential(), &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer{}),
		},
	})
	require.NoError(t, err)
	resp, err := client.Get(context.Background(), resourceGroupName, vmName, nil)
	require.NoError(t, err)
	vm := resp.VirtualMachine

	// the returned VM must satisfy the following conditions
	require.NotNil(t, vm.Name)
	require.Equal(t, vmName, *vm.Name)
	require.NotNil(t, vm.ID)
	require.Contains(t, *vm.ID, resourceGroupName)
}

func (v *fakeVirtualMachinesServer) BeginDelete(ctx context.Context, resourceGroupName string, vmName string, options *armcompute.VirtualMachinesClientBeginDeleteOptions) (resp fake.PollerResponder[armcompute.VirtualMachinesClientDeleteResponse], err fake.ErrorResponder) {
	resp = fake.PollerResponder[armcompute.VirtualMachinesClientDeleteResponse]{}
	resp.AddNonTerminalResponse(nil)
	resp.AddNonTerminalResponse(nil)
	resp.SetTerminalError("VM not existed", http.StatusNotFound)
	return resp, fake.ErrorResponder{}
}

func Test_VirtualMachinesClient_BeginDelete(t *testing.T) {
	// write a fake for VirtualMachinesClient.BeginDelete that satisfies the following requirements

	const (
		vmName            = "virtualmachine1"
		resourceGroupName = "fake-resource-group"
	)

	// TODO: populate pollingErr with the error after polling completes.
	// the fake should include at least one non-terminal response.
	client, err := armcompute.NewVirtualMachinesClient("subscriptionID", azfake.NewTokenCredential(), &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer{}),
		},
	})
	require.NoError(t, err)
	poller, err := client.BeginDelete(context.Background(), resourceGroupName, vmName, nil)
	require.NoError(t, err)
	_, pollingErr := poller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{Frequency: 1 * time.Microsecond})

	// the LRO must terminate in a way to satisfy the following conditions
	require.Error(t, pollingErr)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, pollingErr, &respErr)
}

func (v *fakeVirtualMachinesServer) NewListPager(resourceGroupName string, options *armcompute.VirtualMachinesClientListOptions) (resp fake.PagerResponder[armcompute.VirtualMachinesClientListResponse]) {
	resp = fake.PagerResponder[armcompute.VirtualMachinesClientListResponse]{}
	resp.AddPage(armcompute.VirtualMachinesClientListResponse{
		VirtualMachineListResult: armcompute.VirtualMachineListResult{
			Value: []*armcompute.VirtualMachine{
				{
					Name: to.Ptr("vm1"),
				},
				{
					Name: to.Ptr("vm2"),
				},
				{
					Name: to.Ptr("vm3"),
				},
			},
		},
	}, nil)
	resp.AddError(errors.New("Network issue"))
	resp.AddPage(armcompute.VirtualMachinesClientListResponse{
		VirtualMachineListResult: armcompute.VirtualMachineListResult{
			Value: []*armcompute.VirtualMachine{
				{
					Name: to.Ptr("vm4"),
				},
				{
					Name: to.Ptr("vm5"),
				},
			},
		},
	}, nil)
	return resp
}

func Test_VirtualMachinesClient_NewListPager(t *testing.T) {
	// write a fake for VirtualMachinesClient.NewListPager that satisfies the following requirements

	const (
		resourceGroupName = "fake-resource-group"
	)

	// the fake must return a total of five VMs over two pages.
	// to keep things simple, the returned armcompute.VirtualMachine instances can be empty.
	// while iterating over pages, the fake must return one transient error before the final page

	client, err := armcompute.NewVirtualMachinesClient("subscriptionID", azfake.NewTokenCredential(), &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer{}),
		},
	})
	require.NoError(t, err)

	// TODO: populate vmCount with the number of VMs returned
	var vmCount int

	// TODO: populate pageCount with the number of returned pages
	var pageCount int

	// TODO: populate errCount with the number of transient errors encountered during paging
	var errCount int

	pager := client.NewListPager(resourceGroupName, nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			errCount += 1
		} else {
			pageCount += 1
			for _ = range page.Value {
				vmCount += 1
			}
		}
	}

	// the results must satisfy the following conditions
	require.Equal(t, 5, vmCount)
	require.Equal(t, 2, pageCount)
	require.Equal(t, 1, errCount)
}
