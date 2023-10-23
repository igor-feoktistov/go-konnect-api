package konnect

import (
	"fmt"
	"net/http"
)

type Service struct {
	ConnectTimeout int      `json:"connect_timeout"`
	CreatedAt      int64    `json:"created_at"`
	Enabled	       bool     `json:"enabled"`
	Host           string   `json:"host,omitempty"`
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Path           string   `json:"path,omitempty"`
	Port           int      `json:"port"`
	Protocol       string   `json:"protocol,omitempty"`
	ReadTimeout    int      `json:"read_timeout"`
	Retries        int      `json:"retries"`
	Tags           []string `json:"tags,omitempty"`
	UpdatedAt      int64    `json:"updated_at"`
	WriteTimeout   int      `json:"write_timeout"`
}

type ServiceResponse struct {
	Data   []Service `json:"data"`
	Next   string    `json:"next,omitempty"`
	Offset string    `json:"offset,omitempty"`
}

func (c *Client) ServicesGet(controlPlaneId string, parameters []string) (services []Service, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/services", controlPlaneId)
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		s := ServiceResponse{}
		if _, err = c.Do(req, &s); err != nil {
			return
		}
		services = append(services, s.Data...)
		if len(s.Next) == 0 || len(s.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + s.Offset)
	}
	return
}

func (c *Client) ServiceGetById(controlPlaneId string, serviceId string) (service *Service, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/services/%s", controlPlaneId, serviceId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
			return
	}
	service = &Service{}
	_, err = c.Do(req, service)
	return
}
