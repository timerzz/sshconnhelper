# sshconnhelper
> A tool to create a remote `docker client` via the `ssh client` in `golang.org/x/crypto/ssh`
>
English | [中文](README_ZH.md)

## introduce
Currently, there are two main ways to create a remote `Docker Engine` client through golang:
1. `Docker Engine` opens the remote connection port, and then creates a `docker client` through `client.NewClientWithOpts(client.WithHost("tcp://ip:port"))`.
2. Do ssh password-free on the machine where this node and docker are located, and then create it in the following ways
````golang
helper, _ := connhelper.GetConnectionHelper("ssh://name@ip:port")
   httpClient := &http.Client{
   Transport: &http.Transport{
     DialContext: helper.Dialer,
   },
}
````

Both of the above methods are not very secure. In many cases, the port `2375` will not be opened in the environment, and it is very insecure to do ssh password-free.
In order to solve the above problems, this library can create a `docker client` through `ssh client`, whether it is a `ssh client` created by means of account password or password-free. There is no need to modify the configuration to open port 2375, which is more flexible.

## Install
```powershell
go get github.com/timerzz/sshconnhelper
````

## use

````go
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

````
