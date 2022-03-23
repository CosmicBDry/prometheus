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

type RegisterTask struct {
	Client *reqClient.Client
}

func NewRegisterTask(options *options.Option) *RegisterTask {
	return &RegisterTask{
		Client: reqClient.NewClient(options),
	}
}
func (r *RegisterTask) Run() {
	ticker := time.NewTicker(10 * time.Second)
	logger := Logger.SetLog()
	responsejson := &models.ResponseJson{}
	for {
		r.Client.Options.UUID = options.GetUUID()
		r.Client.Options.Addr = options.GetAddr()
		r.Client.Options.HostName = options.GetHostName()

		event := map[string]interface{}{
			"UUID":     r.Client.Options.UUID,
			"Addr":     r.Client.Options.Addr,
			"HostName": r.Client.Options.HostName,
		}
		path := fmt.Sprintf("%s/agent/register", r.Client.Options.Server)
		//resp, err := r.Client.ClientReq.Post(path, req.BodyJSON(event))
		resp, err := r.Client.ClientReq.Post(path, req.BodyJSON(event))

		json.Unmarshal([]byte(resp.String()), responsejson)
		if responsejson.Code == 200 {
			logger.Info(responsejson.Text)
		} else if responsejson.Code == 201 {
			logger.Info(responsejson.Text)
		} else if responsejson.Code == 400 {
			logger.Error(responsejson.Text)
		} else if err != nil {
			logger.Error(err)
		}
		//fmt.Println(resp.ToString()) ToString同时输出响应body和error
		//fmt.Println(resp.String())
		<-ticker.C
	}
}
