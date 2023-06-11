// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RemoveItemParamsBody remove item params body
//
// swagger:model removeItemParamsBody
type RemoveItemParamsBody struct {

	// product Id
	ProductID string `json:"productId,omitempty"`

	// quantity
	Quantity int32 `json:"quantity,omitempty"`
}

// Validate validates this remove item params body
func (m *RemoveItemParamsBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this remove item params body based on context it is used
func (m *RemoveItemParamsBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RemoveItemParamsBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RemoveItemParamsBody) UnmarshalBinary(b []byte) error {
	var res RemoveItemParamsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
