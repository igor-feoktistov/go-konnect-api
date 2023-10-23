package konnect

import (
	"fmt"
	"net/http"
)

type Headers interface{}

type Healthy struct {
	HttpStatuses []int `json:"http_statuses"`
	Successes    int   `json:"successes"`
	Interval     int   `json:"interval"`
}

type Unhealthy struct {
	HttpStatuses []int `json:"http_statuses"`
	HttpFailures int   `json:"http_failures"`
	Timeouts     int   `json:"timeouts"`
	TcpFailures  int   `json:"tcp_failures"`
	Interval     int   `json:"interval"`
}

type ActiveHealthcheck struct {
	Concurrency            int       `json:"concurrency"`
	Headers                Headers   `json:"headers,omitempty"`
	Healthy                Healthy   `json:"healthy"`
	HttpsVerifyCertificate bool      `json:"https_verify_certificate"`
	HttpPath               string    `json:"http_path,omitempty"`
	HttpsSni               string    `json:"https_sni,omitempty"`
	Timeout                int       `json:"timeout"`
	Type                   string    `json:"type,omitempty"`
	Unhealthy              Unhealthy `json:"unhealthy"`
}

type PassiveHealthcheck struct {
	Healthy   Healthy   `json:"healthy"`
	Type      string    `json:"type,omitempty"`
	Unhealthy Unhealthy `json:"unhealthy"`
}

type Healthchecks struct {
	Active    ActiveHealthcheck  `json:"active"`
	Threshold int                `json:"threshold"`
	Passive   PassiveHealthcheck `json:"passive"`
}

type Upstream struct {
	Algorithm         string          `json:"created_at"`
	ClientCertificate *CertificateRef `json:"client_certificate,omitempty"`
	CreatedAt         int64           `json:"created_at"`
	HashFallback      string          `json:"hash_fallback"`
	HashOn            string          `json:"hash_on"`
	HashOnCookiePath  string          `json:"hash_on_cookie_path"`
	Healthchecks      Healthchecks    `json:"healthchecks"`
	HostHeader        string          `json:""host_header,omitempty"`
	Id                string          `json:"id"`
	Name              string          `json:"name"`
	Slots             int             `json:"slots"`
	Tags              []string        `json:"tags,omitempty"`
	UseSrvName        bool            `json:"use_srv_name"`
}

type UpstreamResponse struct {
	Data   []Upstream `json:"data"`
	Next   string     `json:"next,omitempty"`
	Offset string     `json:"offset,omitempty"`
}

func (c *Client) UpstreamsGet(controlPlaneId string, parameters []string) (upstreams []Upstream, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/upstreams", controlPlaneId)
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		resp := UpstreamResponse{}
		if _, err = c.Do(req, &resp); err != nil {
			return
		}
		upstreams = append(upstreams, resp.Data...)
		if len(resp.Next) == 0 || len(resp.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + resp.Offset)
	}
	return
}

func (c *Client) UpstreamGetById(controlPlaneId string, upstreamId string) (upstream *Upstream, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/upstreams/%s", controlPlaneId, upstreamId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
			return
	}
	upstream = &Upstream{}
	_, err = c.Do(req, upstream)
	return
}
