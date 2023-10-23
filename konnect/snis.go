package konnect

import (
	"fmt"
	"net/http"
)

type CertificateRef struct {
	Id string `json:"id"`
}

type Sni struct {
	Certificate CertificateRef `json:"certificate"`
	CreatedAt   int64          `json:"created_at"`
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	UpdatedAt   int64          `json:"updated_at"`
}

type SniResponse struct {
	Data   []Sni  `json:"data"`
	Next   string `json:"next,omitempty"`
	Offset string `json:"offset,omitempty"`
}

func (c *Client) SnisGet(controlPlaneId string, parameters []string) (snis []Sni, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/snis", controlPlaneId)
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		r := SniResponse{}
		if _, err = c.Do(req, &r); err != nil {
			return
		}
		snis = append(snis, r.Data...)
		if len(r.Next) == 0 || len(r.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + r.Offset)
	}
	return
}

func (c *Client) SniGetById(controlPlaneId string, sniId string) (sni *Sni, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/snis/%s", controlPlaneId, sniId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	sni = &Sni{}
	_, err = c.Do(req, sni)
	return
}
