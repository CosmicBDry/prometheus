package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CosmicBDry/prometheus/configAgent/Logger"
	"github.com/CosmicBDry/prometheus/configAgent/models"
	"github.com/CosmicBDry/prometheus/configAgent/options"
	"github.com/CosmicBDry/prometheus/configAgent/reqClient"
	"github.com/imroc/req"
)

type HeartBeatTask struct {
	Client *reqClient.Client
}

func NewHeartBeatTask(options *options.Option) *HeartBeatTask {

	return &HeartBeatTask{
		Client: reqClient.NewClient(options),
	}

}

func (h *HeartBeatTask) Run() {
	Responsejson := &models.ResponseJson{}
	ticker := time.NewTicker(30 * time.Second)
	logger := Logger.SetLog()
	for {
		path := fmt.Sprintf("%s/agent/heartbeat", h.Client.Options.Server)
		resp, err := h.Client.ClientReq.Get(path, req.Param{"uuid": h.Client.Options.UUID})
		json.Unmarshal([]byte(resp.String()), Responsejson)
		if err != nil {
			logger.Error(err)
		} else if Responsejson.Code == 400 {
			logger.Error(Responsejson.Text)
		} else if Responsejson.Code == 200 {
			logger.Info(Responsejson.Text)
		}

		<-ticker.C
	}

}
