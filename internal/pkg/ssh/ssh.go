package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// Remember to close the client
func InitSshClient(user, password, host string, port uint) (*ssh.Client, error) {
	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // not check host key
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}

	// create ssh client
	addr := fmt.Sprintf("%s:%d", host, port)
	return ssh.Dial("tcp", addr, config)

	//defer sshClient.Close()

}

func Exec(client *ssh.Client, cmd string) (string, error) {
	// create ssh session
	session, err := client.NewSession()
	if err != nil {
		return "", err // create ssh session fail
	}

	defer session.Close()

	// exec cmd
	// combo, err := session.CombinedOutput("whoami; cd /; ls -al;")
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), err // exec cmd fail
	}
	// log.Println("命令输出:", string(combo))
	return "", nil
}
