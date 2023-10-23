package konnect

import (
	"fmt"
	"strconv"
	"net/http"
)

type Page struct {
	Total  int `json:"total"`
	Size   int `json:"size"`
	Number int `json:"number"`
}

type Meta struct {
	Page Page `json:"page"`
}

type ControlPlaneConfig struct {
	ControlPlaneEndpoint string `json:"control_plane_endpoint,omitempty"`
	TelemetryEndpoint    string `json:"telemetry_endpoint,omitempty"`
	ClusterType          string `json:"cluster_type,omitempty"`
}

type ControlPlane struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Config      ControlPlaneConfig `json:"config"`
	CreatedAt   string             `json:"created_at,omitempty"`
	UpdatedAt   string             `json:"updated_at,omitempty"`
}

type ControlPlaneResponse struct {
	Meta Meta `json:"meta"`
	Data []ControlPlane
}

func (c *Client) ControlPlanesGet(parameters []string) (controlPlanes []ControlPlane, err error) {
	var resp *ControlPlaneResponse
	page := 1
	for {
		if resp, err = c.ControlPlanesGetPage(100, page, parameters); err != nil {
			return
		}
		controlPlanes = append(controlPlanes, resp.Data...)
		if len(controlPlanes) == resp.Meta.Page.Total {
			break
		}
		page += 1
	}
	return
}
	
func (c *Client) ControlPlanesGetPage(pageSize int, pageNumber int, parameters []string) (*ControlPlaneResponse, error) {
	var err error
	var req *http.Request
	r := ControlPlaneResponse{}
	path := "/v2/control-planes"
	parameters = append(parameters, "page[size]=" + strconv.Itoa(pageSize))
	parameters = append(parameters, "page[number]=" + strconv.Itoa(pageNumber))
	if req, err = c.NewRequest("GET", path, parameters, nil); err != nil {
		return nil, err
	}
	if _, err = c.Do(req, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) ControlPlaneGetId(controlPlaneId string) (controlPlane *ControlPlane, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s", controlPlaneId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	controlPlane = &ControlPlane{}
	_, err = c.Do(req, controlPlane)
	return
}
