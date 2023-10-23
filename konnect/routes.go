package konnect

import (
	"fmt"
	"net/http"
)

type ServiceRef struct {
	Id string `json:"id"`
}

type Route struct {
	CreatedAt               int64      `json:"created_at"`
	Hosts                   []string   `json:"hosts,omitempty"`
	HttpsRedirectStatusCode int        `json:"https_redirect_status_code"`
	Id                      string     `json:"id"`
	Name                    string     `json:"name"`
	PathHandling            string     `json:"path_handling,omitempty"`
	Paths                   []string   `json:"paths,omitempty"`
	PreserveHost            bool       `json:"preserve_host"`
	Protocols               []string   `json:"protocols,omitempty"`
	RegexPriority           int        `json:"regex_priority"`
	RequestBuffering        bool       `json:"request_buffering"`
	ResponseBuffering       bool       `json:"response_buffering"`
	Service                 ServiceRef `json:"service"`
	StripPath               bool       `json:"strip_path"`
	Tags                    []string   `json:"tags,omitempty"`
	UpdatedAt               int64      `json:"updated_at"`
}

type RouteResponse struct {
	Data   []Route `json:"data"`
	Next   string  `json:"next,omitempty"`
	Offset string  `json:"offset,omitempty"`
}

func (c *Client) RoutesGet(controlPlaneId string, parameters []string) (routes []Route, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/routes", controlPlaneId)
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		r := RouteResponse{}
		if _, err = c.Do(req, &r); err != nil {
			return
		}
		routes = append(routes, r.Data...)
		if len(r.Next) == 0 || len(r.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + r.Offset)
	}
	return
}

func (c *Client) RouteGetById(controlPlaneId string, routeId string) (route *Route, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/routes/%s", controlPlaneId, routeId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	route = &Route{}
	_, err = c.Do(req, route)
	return
}
