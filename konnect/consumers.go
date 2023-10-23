package konnect

import (
	"fmt"
	"net/http"
)

type Consumer struct {
	CreatedAt int64      `json:"created_at"`
	CustomId  string     `json:"custom_id"`
	Id        string     `json:"id"`
	Tags      []string   `json:"tags,omitempty"`
	UpdatedAt int64      `json:"updated_at"`
	UserName  string     `json:""username"`
}

type ConsumerResponse struct {
	Data   []Consumer `json:"data"`
	Next   string     `json:"next,omitempty"`
	Offset string     `json:"offset,omitempty"`
}

func (c *Client) ConsumersGet(controlPlaneId string, parameters []string) (consumers []Consumer, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/consumers", controlPlaneId)
	currentParameters := parameters
	for {
		if req, err = c.NewRequest("GET", path, currentParameters, nil); err != nil {
			return
		}
		r := ConsumerResponse{}
		if _, err = c.Do(req, &r); err != nil {
			return
		}
		consumers = append(consumers, r.Data...)
		if len(r.Next) == 0 || len(r.Offset) == 0 {
			break
		}
		currentParameters = append(parameters, "offset=" + r.Offset)
	}
	return
}

func (c *Client) ConsumerGetById(controlPlaneId string, consumerId string) (consumer *Consumer, err error) {
	var req *http.Request
	path := fmt.Sprintf("/v2/control-planes/%s/core-entities/consumers/%s", controlPlaneId, consumerId)
	if req, err = c.NewRequest("GET", path, []string{}, nil); err != nil {
		return
	}
	consumer = &Consumer{}
	_, err = c.Do(req, consumer)
	return
}
