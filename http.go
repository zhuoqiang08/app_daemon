package main

import (
	"fmt"
	"net/http"
	"strings"
)

//实现对客户端的安装升级功能。
type Http_Service struct {
}

//为虚拟机内部提供服务
func (this *Http_Service) Start() {
	http.HandleFunc("/api/config", this.Service) //设置访问的路由
	err := http.ListenAndServe(":8888", nil)     //设置监听的端口
	if err != nil {
		fmt.Println("(this *Http_Service) Start", err)
	}
}

const (
	WORKDIR = "/opt/smart_stream/agent"
)

func (this *Http_Service) Service(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	op := r.Form.Get("op")
	if len(op) == 0 {
		fmt.Fprintf(w, "error")
		return
	}
	//接收到要求修改网络信息
	if op == "INSTALL" { //安装应用

		content := r.PostForm.Get("content")
		fmt.Println(content)
	}
	if op == "CONFIG" { //配置信息
		content := r.Form.Get("content")
		fmt.Println("接收到配置信息...\n", content)
		if strings.Contains(content, "agent") &&
			strings.Contains(content, "nsq_server") &&
			strings.Contains(content, "config_path") &&
			strings.Contains(content, "local_ip") {
			WriteFile("./app.conf", content)
		}
	}
	//收到升级信息

	if op == "cmd" {
		content := r.Form.Get("content")
		ExecCmd(string(content))
	}
	fmt.Fprintf(w, "success") //输出到客户端的信息
}

func (this *Http_Service) update_agent(url string) {

	ExecCmd("mkdir " + WORKDIR + "/tmp")
	ExecCmd("rm -rf " + WORKDIR + "/tmp/stream_agent")
	ExecCmd(fmt.Sprintf("wget -O "+WORKDIR+"/tmp/stream_agent %s > /dev/null 2>&1", url))
	ExecCmd("chmod +x " + WORKDIR + "/tmp/stream_agent")
	ExecCmd("pkill stream_agent")
	ExecCmd("/bin/cp -rf " + WORKDIR + "/tmp/stream_agent " + WORKDIR + "/")
	ExecCmd("chmod +x " + WORKDIR + "/stream_agent")
	WriteFile("./run.sh", "./stream_agent &")
	ExecCmd("chmod +x " + WORKDIR + "/run.sh")
	go ExecCmd("./run.sh")

}
func (this *Http_Service) Install_nginx(url string) {

	ExecCmd("pkill nginx")
	ExecCmd("rm  -rf /root/openresty.tar.gz")
	ExecCmd(fmt.Sprintf("wget -o /root/openresty.tar.gz %s > /dev/null 2>&1", url))
	ExecCmd("cp -rf /opt/openresty/nginx/sbin/nginx /sbin/")
	ExecCmd("tar -zxvf /root/openresty.tar.gz -C /")
	go ExecCmd("/opt/openresty/nginx/sbin/nginx &")

}
