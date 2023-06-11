// Code generated by go-swagger; DO NOT EDIT.

package product

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"eda-in-golang/stores/storesclient/models"
)

// GetStoreProductReader is a Reader for the GetStoreProduct structure.
type GetStoreProductReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetStoreProductReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetStoreProductOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetStoreProductDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetStoreProductOK creates a GetStoreProductOK with default headers values
func NewGetStoreProductOK() *GetStoreProductOK {
	return &GetStoreProductOK{}
}

/* GetStoreProductOK describes a response with status code 200, with default header values.

A successful response.
*/
type GetStoreProductOK struct {
	Payload *models.StorespbGetCatalogResponse
}

func (o *GetStoreProductOK) Error() string {
	return fmt.Sprintf("[GET /api/stores/{storeId}/products][%d] getStoreProductOK  %+v", 200, o.Payload)
}
func (o *GetStoreProductOK) GetPayload() *models.StorespbGetCatalogResponse {
	return o.Payload
}

func (o *GetStoreProductOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StorespbGetCatalogResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetStoreProductDefault creates a GetStoreProductDefault with default headers values
func NewGetStoreProductDefault(code int) *GetStoreProductDefault {
	return &GetStoreProductDefault{
		_statusCode: code,
	}
}

/* GetStoreProductDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type GetStoreProductDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// Code gets the status code for the get store product default response
func (o *GetStoreProductDefault) Code() int {
	return o._statusCode
}

func (o *GetStoreProductDefault) Error() string {
	return fmt.Sprintf("[GET /api/stores/{storeId}/products][%d] getStoreProduct default  %+v", o._statusCode, o.Payload)
}
func (o *GetStoreProductDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *GetStoreProductDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}