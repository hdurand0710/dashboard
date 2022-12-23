// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PublicGCPCloudSpec PublicGCPCloudSpec is a public counterpart of apiv1.GCPCloudSpec.
//
// swagger:model PublicGCPCloudSpec
type PublicGCPCloudSpec struct {

	// node ports allowed IP ranges
	NodePortsAllowedIPRanges *NetworkRanges `json:"nodePortsAllowedIPRanges,omitempty"`
}

// Validate validates this public g c p cloud spec
func (m *PublicGCPCloudSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateNodePortsAllowedIPRanges(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicGCPCloudSpec) validateNodePortsAllowedIPRanges(formats strfmt.Registry) error {
	if swag.IsZero(m.NodePortsAllowedIPRanges) { // not required
		return nil
	}

	if m.NodePortsAllowedIPRanges != nil {
		if err := m.NodePortsAllowedIPRanges.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("nodePortsAllowedIPRanges")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("nodePortsAllowedIPRanges")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this public g c p cloud spec based on the context it is used
func (m *PublicGCPCloudSpec) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateNodePortsAllowedIPRanges(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicGCPCloudSpec) contextValidateNodePortsAllowedIPRanges(ctx context.Context, formats strfmt.Registry) error {

	if m.NodePortsAllowedIPRanges != nil {
		if err := m.NodePortsAllowedIPRanges.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("nodePortsAllowedIPRanges")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("nodePortsAllowedIPRanges")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PublicGCPCloudSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PublicGCPCloudSpec) UnmarshalBinary(b []byte) error {
	var res PublicGCPCloudSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}