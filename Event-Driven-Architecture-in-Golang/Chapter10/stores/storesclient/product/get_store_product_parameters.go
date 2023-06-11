// Code generated by go-swagger; DO NOT EDIT.

package product

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetStoreProductParams creates a new GetStoreProductParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetStoreProductParams() *GetStoreProductParams {
	return &GetStoreProductParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetStoreProductParamsWithTimeout creates a new GetStoreProductParams object
// with the ability to set a timeout on a request.
func NewGetStoreProductParamsWithTimeout(timeout time.Duration) *GetStoreProductParams {
	return &GetStoreProductParams{
		timeout: timeout,
	}
}

// NewGetStoreProductParamsWithContext creates a new GetStoreProductParams object
// with the ability to set a context for a request.
func NewGetStoreProductParamsWithContext(ctx context.Context) *GetStoreProductParams {
	return &GetStoreProductParams{
		Context: ctx,
	}
}

// NewGetStoreProductParamsWithHTTPClient creates a new GetStoreProductParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetStoreProductParamsWithHTTPClient(client *http.Client) *GetStoreProductParams {
	return &GetStoreProductParams{
		HTTPClient: client,
	}
}

/* GetStoreProductParams contains all the parameters to send to the API endpoint
   for the get store product operation.

   Typically these are written to a http.Request.
*/
type GetStoreProductParams struct {

	// StoreID.
	StoreID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get store product params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetStoreProductParams) WithDefaults() *GetStoreProductParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get store product params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetStoreProductParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get store product params
func (o *GetStoreProductParams) WithTimeout(timeout time.Duration) *GetStoreProductParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get store product params
func (o *GetStoreProductParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get store product params
func (o *GetStoreProductParams) WithContext(ctx context.Context) *GetStoreProductParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get store product params
func (o *GetStoreProductParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get store product params
func (o *GetStoreProductParams) WithHTTPClient(client *http.Client) *GetStoreProductParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get store product params
func (o *GetStoreProductParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithStoreID adds the storeID to the get store product params
func (o *GetStoreProductParams) WithStoreID(storeID string) *GetStoreProductParams {
	o.SetStoreID(storeID)
	return o
}

// SetStoreID adds the storeId to the get store product params
func (o *GetStoreProductParams) SetStoreID(storeID string) {
	o.StoreID = storeID
}

// WriteToRequest writes these params to a swagger request
func (o *GetStoreProductParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param storeId
	if err := r.SetPathParam("storeId", o.StoreID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
