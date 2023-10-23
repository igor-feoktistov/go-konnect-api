package konnect

import (
	"fmt"
	"net/http"
)

type RouteRef struct {
	Id string `json:"id"`
}

type Config interface{}

type Plugin struct {
	Config    Config     `json:"config"`
	CreatedAt int64      `json:"created_at"`
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Protocols []string   `json:"protocols,omitempty"`
	Route     ServiceRef `json:"service,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	UpdatedAt int64      `json:"updated_at"`
}

type PluginResponse struct {
	Data   []Plugin `json:"data"`
	Next   string   `json:"next,omitempty"`
	Offset string   `json:"offset,omitempty"`
}

func (c *Client) PluginsGet(controlPlaneId string, routeId string, parameters []string) (plugins []Plugin, err error) {
	var req *http.Request
	var path string
	if len(routeId) > 0 {
		path = fmt.Sprintf("/v2/control-planes/%s/core-entities/routes/%s/plugins", controlPlaneId, routeId)
	} else {
		path = fmt.Sprintf("/v2/control-planes/%s/core-entities/plugins", controlPlaneId)
	}
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		r := PluginResponse{}
		if _, err = c.Do(req, &r); err != nil {
			return
		}
		plugins = append(plugins, r.Data...)
		if len(r.Next) == 0 || len(r.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + r.Offset)
	}
	return
}

func (c *Client) PluginGetById(controlPlaneId string, routeId string, pluginId string) (plugin *Plugin, err error) {
	var req *http.Request
	var path string
	if len(routeId) > 0 {
		path = fmt.Sprintf("/v2/control-planes/%s/core-entities/routes/%s/plugins/%s", controlPlaneId, routeId, pluginId)
	} else {
		path = fmt.Sprintf("/v2/control-planes/%s/core-entities/plugins/%s", controlPlaneId, pluginId)
	}
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	plugin = &Plugin{}
	_, err = c.Do(req, plugin)
	return
}
