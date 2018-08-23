package main

//linux
import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

func ExecCmd(cmdline string) (string, error) {
	var cmd *exec.Cmd
	log.Println("ExecCmd:", string(cmdline))
	cmd = exec.Command("/bin/bash", "-c", cmdline)
	out, err := cmd.CombinedOutput()
	log.Println("out:", string(out))
	return string(out), err
}

func ExecCmd_File(cmdline string) (string, error) {
	WriteFile("./tmp.sh", cmdline)
	ExecCmd("chmod +x ./tmp.sh")
	return ExecCmd("./tmp.sh")
}

func ExecCmd_Nsq(cmdline string, msg_id string, nsq *Nsq_Service) (string, error) {
	out, err := ExecCmd(cmdline)
	ret := MQ_Ret{}
	if err != nil {
		ret.Error = err.Error()
	}
	ret.Msg = out
	ret.Id = msg_id
	ret.success = true
	if err == nil {
		ret.success = false
	}
	info, _ := json.Marshal(&ret)
	nsq.Publish(RET_TOPIC, string(info))
	return out, err
}
func ExecCmd_File_Nsq(cmdline string, msg_id string, nsq *Nsq_Service) (string, error) {
	WriteFile("./tmp.sh", cmdline)
	ExecCmd("chmod +x ./tmp.sh")
	out, err := ExecCmd("./tmp.sh")
	ret := MQ_Ret{}
	if err != nil {
		ret.Error = err.Error()
	}
	ret.Msg = out
	ret.Id = msg_id
	ret.success = true
	if err == nil {
		ret.success = false
	}
	info, _ := json.Marshal(&ret)
	nsq.Publish(RET_TOPIC, string(info))
	return out, err
}

func WriteFile(filepath, content string) {
	file, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		log.Println(err)
	}
	file.Close()
}
