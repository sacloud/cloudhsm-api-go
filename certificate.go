// Copyright 2025- The sacloud/cloudhsm-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cloudhsm

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
)

type ClientAPI interface {
	List(ctx context.Context) ([]v1.CloudHSMClient, error)
	Create(ctx context.Context, request CloudHSMClientCreateParams) (*v1.CloudHSMClient, error)
	Read(ctx context.Context, id string) (*v1.CloudHSMClient, error)
	Update(ctx context.Context, id string, params CloudHSMClientUpdateParams) (*v1.CloudHSMClient, error)
	Delete(ctx context.Context, id string) error
}

var _ ClientAPI = (*ClientOp)(nil)

type ClientOp struct {
	client *v1.Client
	hsm    *v1.CloudHSM
}

func NewClientOp(client *v1.Client, hsm *v1.CloudHSM) (ClientAPI, error) {
	if hsm.GetAvailability() == v1.AvailabilityEnumAvailable {
		return &ClientOp{
			client: client,
			hsm:    hsm,
		}, nil
	}
	return nil, errors.New("CloudHSM unavailable")
}

func (op *ClientOp) List(ctx context.Context) ([]v1.CloudHSMClient, error) {
	resp, err := op.client.CloudhsmCloudhsmsClientsList(
		ctx,
		v1.CloudhsmCloudhsmsClientsListParams{
			CloudhsmResourceID: op.hsm.GetID(),
		},
	)

	if err == nil {
		return resp.GetClients(), nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSM.Read", 0, err)
	} else {
		return nil, NewAPIError("CloudHSM.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

type CloudHSMClientCreateParams struct {
	Name        string
	Certificate string
}

func (op *ClientOp) Create(ctx context.Context, p CloudHSMClientCreateParams) (*v1.CloudHSMClient, error) {
	resp, err := op.client.CloudhsmCloudhsmsClientsCreate(
		ctx,
		&v1.WrappedCreateCloudHSMClient{
			Client: v1.CreateCloudHSMClient{
				Name:         p.Name,
				Certificate:  p.Certificate,
				Availability: v1.AvailabilityEnumPrecreate,
			},
		},
		v1.CloudhsmCloudhsmsClientsCreateParams{
			CloudhsmResourceID: op.hsm.GetID(),
		},
	)

	if err == nil {
		c := resp.GetClient()
		// Convert v1.CreateCloudHSMClient to v1.CloudHSMClient
		client := v1.CloudHSMClient{
			ID:           c.ID,
			CreatedAt:    c.CreatedAt,
			ModifiedAt:   c.ModifiedAt,
			Availability: c.Availability,
			Name:         c.Name,
			Certificate:  c.Certificate,
		}
		return &client, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSMClient.Create", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("CloudHSMClient.Create", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("CloudHSMClient.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *ClientOp) Read(ctx context.Context, id string) (*v1.CloudHSMClient, error) {
	resp, err := op.client.CloudhsmCloudhsmsClientsRetrieve(
		ctx,
		v1.CloudhsmCloudhsmsClientsRetrieveParams{
			CloudhsmResourceID: op.hsm.GetID(),
			ID:                 id,
		},
	)

	if err == nil {
		client := resp.GetClient()
		return &client, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSMClient.Read", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return nil, NewAPIError("CloudHSMClient.Read", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return nil, NewAPIError("CloudHSMClient.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

type CloudHSMClientUpdateParams struct {
	Name string
}

func (op *ClientOp) Update(ctx context.Context, id string, p CloudHSMClientUpdateParams) (*v1.CloudHSMClient, error) {
	resp, err := op.client.CloudhsmCloudhsmsClientsUpdate(
		ctx,
		&v1.WrappedCloudHSMClient{
			Client: v1.CloudHSMClient{
				Name: p.Name,

				// This cannot be updated but zero is invalid...
				Availability: v1.AvailabilityEnumAvailable,
			},
		},
		v1.CloudhsmCloudhsmsClientsUpdateParams{
			CloudhsmResourceID: op.hsm.GetID(),
			ID:                 id,
		},
	)

	if err == nil {
		client := resp.GetClient()
		return &client, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSMClient.Update", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("CloudHSMClient.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("CloudHSMClient.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *ClientOp) Delete(ctx context.Context, id string) error {
	err := op.client.CloudhsmCloudhsmsClientsDestroy(
		ctx,
		v1.CloudhsmCloudhsmsClientsDestroyParams{
			CloudhsmResourceID: op.hsm.GetID(),
			ID:                 id,
		},
	)

	if err == nil {
		return nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return NewAPIError("CloudHSMClient.Delete", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return NewAPIError("CloudHSMClient.Delete", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return NewAPIError("CloudHSMClient.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}
