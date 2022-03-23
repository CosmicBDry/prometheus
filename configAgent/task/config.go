package task

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"os"

	"github.com/CosmicBDry/prometheus/configAgent/Logger"
	"github.com/CosmicBDry/prometheus/configAgent/models"
	"github.com/CosmicBDry/prometheus/configAgent/options"
	"github.com/CosmicBDry/prometheus/configAgent/reqClient"
	"github.com/CosmicBDry/prometheus/configAgent/utils"
	"github.com/imroc/req"
)

type CongfigTask struct {
	Client *reqClient.Client
}

func NewConfigTask(options *options.Option) *CongfigTask {
	return &CongfigTask{
		Client: reqClient.NewClient(options),
	}

}

func (c *CongfigTask) Run() {
	logger := Logger.SetLog()
	ticker := time.NewTicker(10 * time.Second)
	for {
		c.Client.Options.ConfigVersion = options.GetConfigVersion()
		c.Client.Options.UUID = options.GetUUID()
		configModels := models.ResponseJson{}
		path := fmt.Sprintf("%s/agent/getconfig", c.Client.Options.Server)
		resp, _ := c.Client.ClientReq.Get(path, req.Param{
			"uuid":          c.Client.Options.UUID,
			"configversion": c.Client.Options.ConfigVersion,
		})
		json.Unmarshal([]byte(resp.String()), &configModels)
		//判断获取的config版本是否大于本地版本，若大于则更新本地配置文件
		result := c.Client.Options.UpdateConfigVersion(configModels.ConfigVersion)
		if result {
			//ctx, _ := yaml.Marshal(&configModels.Config)
			//fmt.Println(configModels.Config)
			Time1 := time.Now().Unix()
			Time2 := time.Now().Format("2006-01-02_15-04")
			Path := fmt.Sprintf("%s/%s-%d%s", filepath.Dir(c.Client.Options.Conf), "promtheus", Time1, ".yml")
			err := utils.WriteFile(configModels.Config, Path)
			if err != nil {
				logger.Info(err)
			}
			Result, err := utils.CmdRun(fmt.Sprintf("%s %s %s", c.Client.Options.Promtool, "check config", Path))

			if err == nil {
				logger.Info(Result)
				err = os.Rename(c.Client.Options.Conf, fmt.Sprintf("%s-bak-%s", c.Client.Options.Conf, Time2))
				if err == nil {
					logger.Info("configuration backed up!: " + fmt.Sprintf("%s-bak-%s", c.Client.Options.Conf, Time2))
				} else {
					logger.Error(err)
				}
				err = os.Rename(Path, c.Client.Options.Conf)
				if err == nil {
					logger.Info("configuration updated success!: " + c.Client.Options.Conf)
				} else {
					logger.Error(err)
				}

				//重载配置
				_, err = utils.CmdRun("kill -SIGHUP `ps -ef | grep prometheus| grep -v grep|awk '{print $2}'`")
				if err != nil {
					logger.Error("config reload failed!")
				} else {
					logger.Info("config reload success!")
				}
			} else {
				logger.Error(Result)
			}
		}

		//fmt.Println(configModels.Config.ScrapeConfigs[3].BasicAuth)
		<-ticker.C
	}

}
