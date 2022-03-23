package options

import (
	"net"
	"os"
	"strings"

	"github.com/CosmicBDry/prometheus/configAgent/utils"
	"github.com/google/uuid"
)

type Option struct {
	Server        string
	UUID          string
	Addr          string
	HostName      string
	ConfigVersion string
	Promtool      string
	Conf          string
}

//获取本机的uuid
func GetUUID() string {
	UUID := utils.ReadFile("./agentinfo/UUID.info")
	if UUID != "" {
		return UUID
	}
	UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	utils.WriteFile(UUID, "./agentinfo/UUID.info")
	return UUID
}

//获取本地配置文件版本号
func GetConfigVersion() string {
	Version := utils.ReadFile("./agentinfo/ConfigVersion.info")
	if Version != "" {
		return Version
	}
	Version = "0"
	return Version
}

//更新本地配置文件版本方法
func (option *Option) UpdateConfigVersion(version string) bool {

	if option.ConfigVersion < version {

		utils.WriteFile(version, "./agentinfo/ConfigVersion.info")
		return true
	}
	return false

}

func GetHostName() string {
	hostname, _ := os.Hostname()
	return hostname
}

func GetAddr() string {
	IPADDR, _ := net.InterfaceAddrs()
	agentaddr := ""
	for _, addr := range IPADDR {
		if strings.Index(addr.String(), ":") > 0 {
			continue
		}
		nodes := strings.SplitN(addr.String(), "/", 2)
		if len(nodes) != 2 {
			continue
		}
		agentaddr = nodes[0]
	}
	return agentaddr
}

func NewOption(server, promtool, conf string) *Option {
	return &Option{
		Server:        server,
		UUID:          GetUUID(),
		Addr:          GetAddr(),
		HostName:      GetHostName(),
		ConfigVersion: GetConfigVersion(),
		Promtool:      promtool,
		Conf:          conf,
	}
}
