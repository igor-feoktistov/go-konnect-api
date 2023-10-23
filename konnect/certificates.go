package konnect

import (
	"fmt"
	"net/http"
)

type CertMetadata struct {
	DnsNames  []string `json:"dns_names,omitempty"`
	Expiry    string   `json:"expiry"`
	Issuer	  string   `json:"issuer"`
	KeyUsages []string `json:"key_usages"`
	SanNames  []string `json:"san_names"`
	Snis      []string `json:"snis"`
	Subject   string   `json:"subject"`
}

type Certificate struct {
	Cert      string       `json:"cert"`
	CreatedAt int64        `json:"created_at"`
	Id        string       `json:"id"`
	Key       string       `json:"key"`
	Metadata  CertMetadata `json:"metadata"`
	Tags      []string     `json:"tags,omitempty"`
	UpdatedAt int64        `json:"updated_at"`
}

type CertificateResponse struct {
	Data   []Certificate `json:"data"`
	Next   string        `json:"next,omitempty"`
	Offset string        `json:"offset,omitempty"`
}

func (c *Client) CertificatesGet(controlPlaneId string, parameters []string) (certificates []Certificate, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/certificates", controlPlaneId)
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		r := CertificateResponse{}
		if _, err = c.Do(req, &r); err != nil {
			return
		}
		certificates = append(certificates, r.Data...)
		if len(r.Next) == 0 || len(r.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + r.Offset)
	}
	return
}

func (c *Client) CertificateGetById(controlPlaneId string, certificateId string) (certificate *Certificate, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/certificates/%s", controlPlaneId, certificateId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	certificate = &Certificate{}
	_, err = c.Do(req, certificate)
	return
}
