package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-registry/client"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/address"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/google/network"
	"github.com/frodenas/bosh-google-cpi/google/target_pool"
)

type DeleteVM struct {
	vmService         ginstance.GoogleInstanceService
	addressService    gaddress.GoogleAddressService
	networkService    gnetwork.GoogleNetworkService
	targetPoolService gtargetpool.GoogleTargetPoolService
	registryClient    registry.Client
}

func NewDeleteVM(
	vmService ginstance.GoogleInstanceService,
	addressService gaddress.GoogleAddressService,
	networkService gnetwork.GoogleNetworkService,
	targetPoolService gtargetpool.GoogleTargetPoolService,
	registryClient registry.Client,
) DeleteVM {
	return DeleteVM{
		vmService:         vmService,
		addressService:    addressService,
		networkService:    networkService,
		targetPoolService: targetPoolService,
		registryClient:    registryClient,
	}
}

func (dv DeleteVM) Run(vmCID VMCID) (interface{}, error) {
	// Delete VM networks
	var networks Networks
	vmNetworks := networks.AsGoogleInstanceNetworks()
	instanceNetworks := ginstance.NewGoogleInstanceNetworks(vmNetworks, dv.addressService, dv.networkService, dv.targetPoolService)

	err := dv.vmService.DeleteNetworkConfiguration(string(vmCID), instanceNetworks)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM
	err = dv.vmService.Delete(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM agent settings
	err = dv.registryClient.Delete(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	return nil, nil
}
