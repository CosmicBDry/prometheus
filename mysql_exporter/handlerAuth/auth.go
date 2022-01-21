package handlerAuth

import (
	"encoding/base64" //base64解码、编码等
	"net/http"
	"strings"

	"github.com/siruspen/logrus"
	"golang.org/x/crypto/bcrypt" ///bcrypt密码加密的hash运算
)

type AuthSecrets map[string]string //定义一个map类型，用于外部设置访问的user和password，是基于Prometheus的Basic认证类型

//定义一个Auth函数，参数类型和返回值均为http.Handler接口类型---------------------------------------------------------------------->
func Auth(handler http.Handler, secrets AuthSecrets, logger *logrus.Logger) http.Handler {

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		secret := request.Header.Get("Authorization")

		if !IsAuth(secret, secrets) {
			logger.Error("basic_auth认证失败，请提供正确的用户名、密码或有效的api请求token")
			response.Header().Set("WWW-Authenticate", `Basic realm=""`) //响应头设置Basic认证类型：Basic base64(user:password)，并弹出输入框
			response.WriteHeader(401)                                   //状态码必须设置成401，才能弹出上面的Basic弹出框
			return
		}
		logger.Info(request.Method, " ", "http://", request.Host, request.URL)
		logger.Info("User-Agent: ", request.Header["User-Agent"], " ", "ClientAddr: ", "[", request.RemoteAddr, "]")
		handler.ServeHTTP(response, request)

	})
}

//验证客户端用户名及密码的Basic认证---------------------------------------------------------------------->
//判断客户端发送的secret是否等于设置的secrets
func IsAuth(secret string, secrets AuthSecrets) bool {

	if secrets == nil { //若管理员不设置用户名及密码认证则直接返回true，则任何人均可访问，不设限制
		return true
	}

	//通过Fields函数，以空格为分割符，将secret中的Basic标记和认证base64编码分成两个部分，均保存在[]string
	clientAuth := strings.Fields(secret)
	if len(clientAuth) != 2 {
		return false
	}

	//clientAuth[1]为加密的base64编码，将其解码成明文
	plaintext, err := base64.StdEncoding.DecodeString(clientAuth[1])

	if err != nil {
		return false
	}

	//SplitN以冒号为分割符，分割成两个部分，分别是user和password
	clientAuth = strings.SplitN(string(plaintext), ":", 2)
	if len(clientAuth) != 2 {
		return false
	}

	//判断用户名是否存在于secrets映射中，若存在则将secrets中对应的用户密码的hash值赋值给Hashpass
	Hashpass, ok := secrets[clientAuth[0]]

	//判断客户端输入的用户名存在并且密码校验正确，同时满足两个条件方可返回true，允许访问exporter资源
	return ok && bcrypt.CompareHashAndPassword([]byte(Hashpass), []byte(clientAuth[1])) == nil

}
