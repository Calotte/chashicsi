/*
Copyright 2024 chashi.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package csi

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

var (
	controllerCaps = []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	}
)

type controllerService struct {
	csi.UnimplementedControllerServer
}

var _ csi.ControllerServer = &controllerService{}

func newControllerService() controllerService {
	return controllerService{}
}

func (d *controllerService) CreateVolume(ctx context.Context, request *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if len(request.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Name cannot be empty")
	}
	if request.VolumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capabilities cannot be empty")
	}

	requiredCap := request.CapacityRange.GetRequiredBytes()
	klog.Infof("Creating volume: %s, capacity: %dbytes", request.Name, requiredCap)

	// Create container logic is here

	klog.Infof("Created volume: %s ok", request.Name)
	volCtx := make(map[string]string)
	for k, v := range request.Parameters {
		volCtx[k] = v
	}

	volCtx["subPath"] = request.Name
	volume := csi.Volume{
		VolumeId:      request.Name,
		CapacityBytes: requiredCap,
		VolumeContext: volCtx,
	}

	return &csi.CreateVolumeResponse{Volume: &volume}, nil
}

func (d *controllerService) DeleteVolume(ctx context.Context, request *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	klog.Infof("Deleting volume: %s", request.VolumeId)

	// Delete container logic is here

	klog.Infof("Deleted volume: %s ok", request.VolumeId)
	return &csi.DeleteVolumeResponse{}, nil
}

func (d *controllerService) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	var caps []*csi.ControllerServiceCapability
	for _, cap := range controllerCaps {
		c := &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		}
		caps = append(caps, c)
	}
	return &csi.ControllerGetCapabilitiesResponse{Capabilities: caps}, nil
}

func (d *controllerService) ControllerPublishVolume(ctx context.Context, request *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	volumeID := request.VolumeId
	nodeID := request.NodeId
	isAttached, err := d.isVolumeAttached(volumeID, nodeID)
	if err != nil {
		klog.Errorf("failed to check volume attachment: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check volume attachment: %v", err)
	}
	if isAttached {
		klog.Infof("Volume %s is already attached to node %s", volumeID, nodeID)
		return &csi.ControllerPublishVolumeResponse{}, nil
	}
	err = d.attachVolumeToNode(volumeID, nodeID)
	if err != nil {
		klog.Errorf("failed to attach volume %s to node %s: %v", volumeID, nodeID, err)
		return nil, status.Errorf(codes.Internal, "failed to attach volume %s to node %s: %v", volumeID, nodeID, err)
	}
	return &csi.ControllerPublishVolumeResponse{}, nil
}

func (d *controllerService) ControllerUnpublishVolume(ctx context.Context, request *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) ValidateVolumeCapabilities(ctx context.Context, request *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) ListVolumes(ctx context.Context, request *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) GetCapacity(ctx context.Context, request *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) CreateSnapshot(ctx context.Context, request *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) DeleteSnapshot(ctx context.Context, request *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) ListSnapshots(ctx context.Context, request *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) ControllerExpandVolume(ctx context.Context, request *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) ControllerModifyVolume(ctx context.Context, request *csi.ControllerModifyVolumeRequest) (*csi.ControllerModifyVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (d *controllerService) isVolumeAttached(volumeID, nodeID string) (bool, error) {
	klog.Infof("Checking if volume %s is attached to node %s", volumeID, nodeID)
	return false, nil
}

func (d *controllerService) attachVolumeToNode(volumeID, nodeID string) error {
	klog.Infof("Attaching volume %s to node %s", volumeID, nodeID)
	return nil
}
