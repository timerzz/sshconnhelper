# sshconnhelper
> 一个通过`golang.org/x/crypto/ssh`中的`ssh client`创建远程`docker client`的工具
> 
[English](README.md) | 中文

## 介绍
目前通过golang创建远程`Docker Engine`的客户端主要有两种方法：
1. `Docker Engine`开启远程连接端口，然后通过`client.NewClientWithOpts(client.WithHost("tcp://ip:port"))`来创建一个`docker client`。
2. 在本节点和docker所在的机器做ssh免密，然后通过以下方式创建
```golang
helper, _ := connhelper.GetConnectionHelper("ssh://name@ip:port")
   httpClient := &http.Client{
   Transport: &http.Transport{
     DialContext: helper.Dialer,
   },
}
```  

上面两种方式都不是很安全。很多时候环境中不会开通`2375`这个端口，而做ssh免密也十分不安全。  
为了解决上述问题，这个库可以通过`ssh client`创建一个`docker client`，无论是通过账号密码还是免密等方式创建的`ssh client`都可以。不需要修改配置开放2375端口，更加灵活。

## 安装
```powershell
go get github.com/timerzz/sshconnhelper
```

## 使用

```go
package main

import (
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
	"github.com/timerzz/sshconnhelper"
)

func NewClientBySSHClient(cli *ssh.Client) (*client.Client, error) {
	helper := sshconnhelper.GetConnectionHelperBySshClient(cli)
	return client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
		client.WithDialContext(helper.Dialer),
	)
}

```