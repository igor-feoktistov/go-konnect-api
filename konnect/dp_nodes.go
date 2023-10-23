package konnect

import (
	"fmt"
	"net/http"
)

type DpNode struct {
	CompatibilityStatus struct {
		State       string   `json:"state"`
	}                            `json:"compatibility_status"`
	ConfigHash          string   `json:"config_hash"`
	CreatedAt           int64    `json:"created_at"`
	Hostname            string   `json:"hostname"`
	Id                  string   `json:"id"`
	LastPing            int64    `json:"last_ping"`
	Status              string   `json:"status,omitempty"`
	Type                string   `json:"type"`
	UpdatedAt           int64    `json:"updated_at"`
	Version             string   `json:"version"`
}

type DpNodeResponse struct {
	Items   []DpNode       `json:"items"`
	Page struct {
		TotalCount int `json:"total_count"`
	}                      `json:"page"`
}

func (c *Client) DpNodesGet(controlPlaneId string, parameters []string) (nodes []DpNode, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/nodes", controlPlaneId)
	if req, err = c.NewRequest("GET", path, parameters, nil); err != nil {
		return
	}
	r := DpNodeResponse{}
	if _, err = c.Do(req, &r); err == nil {
		nodes = append(nodes, r.Items...)
	}
	return
}

func (c *Client) DpNodeGetById(controlPlaneId string, dpNodeId string) (dpNode *DpNode, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/nodes/%s", controlPlaneId, dpNodeId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	dpNode = &DpNode{}
	_, err = c.Do(req, dpNode)
	return
}
