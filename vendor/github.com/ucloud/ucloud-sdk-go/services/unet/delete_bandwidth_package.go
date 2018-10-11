//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UNet DeleteBandwidthPackage

package unet

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DeleteBandwidthPackageRequest is request schema for DeleteBandwidthPackage action
type DeleteBandwidthPackageRequest struct {
	request.CommonBase

	// 带宽包资源ID
	BandwidthPackageId *string `required:"true"`
}

// DeleteBandwidthPackageResponse is response schema for DeleteBandwidthPackage action
type DeleteBandwidthPackageResponse struct {
	response.CommonBase
}

// NewDeleteBandwidthPackageRequest will create request of DeleteBandwidthPackage action.
func (c *UNetClient) NewDeleteBandwidthPackageRequest() *DeleteBandwidthPackageRequest {
	req := &DeleteBandwidthPackageRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DeleteBandwidthPackage - 删除弹性IP上已附加带宽包
func (c *UNetClient) DeleteBandwidthPackage(req *DeleteBandwidthPackageRequest) (*DeleteBandwidthPackageResponse, error) {
	var err error
	var res DeleteBandwidthPackageResponse

	err = c.client.InvokeAction("DeleteBandwidthPackage", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}