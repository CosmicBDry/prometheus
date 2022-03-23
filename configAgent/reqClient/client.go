package reqClient

import (
	"github.com/CosmicBDry/prometheus/configAgent/options"
	"github.com/imroc/req"
)

type Client struct {
	Options   *options.Option
	ClientReq *req.Req
}

func NewClient(options *options.Option) *Client {
	return &Client{
		Options:   options,
		ClientReq: req.New(),
	}
}
