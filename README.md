gopush-cluster
==============
Golang实现的推送集群服务。

* 安装

## 安装
### 一、搭建zookeeper
1.下载[zookeeper](http://www.apache.org/dyn/closer.cgi/zookeeper/)，推荐下载3.4.5版本

2.解压包（这里解压到 /data/programfiles/zookeeper-3.4.5下面 ）
```sh
$ mkdir -p /data/programfiles
$ cp ./zookeeper-3.4.5.tar.gz /data/programfiles
$ cd /data/programfiles/
$ tar -xvf zookeeper-3.4.5.tar.gz -C ./
```
3.编译及安装
``` sh
$ cd zookeeper-3.4.5/src/c
$ ./configure
$ make && make install
```
4.启动zookeeper(zookeeper配置在这里不做详细介绍)
```sh
$ cd /data/programfiles/zookeeper-3.4.5/bin
$ nohup ./zkServer.sh start &
```
### 二、搭建redis
```sh
$ cd /data/programfiles
$ wget https://redis.googlecode.com/files/redis-2.6.4.tar.gz
$ tar -xvf redis-2.6.4.tar.gz -C ./
$ cd redis-2.6.4
$ make
$ make test
$ make install
$ mkdir /etc/redis
$ cp /data/programfiles/redis-2.6.4/redis.conf /etc/redis/
$ cp /data/programfiles/redis-2.6.4/redis-server /etc/init.d/redis-server
$ /etc/init.d/redis-server /etc/redis/redis.conf
```
### 三、安装git工具（如果已安装则可跳过此步）
参考：[git](http://git-scm.com/download/linux)
### 四、搭建golang环境
1.下载源码(根据自己的系统下载对应的安装包)
```sh
$ wget https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
$ tar -xvf go1.2.linux-amd64.tar.gz
$ cp -R go /usr/local/
```
2.配置GO环境变量
(这里我加在/etc/profile)
```sh
$ vim etc/profile
# 将以下环境变量添加到profile最后面
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/data/app/go
```
### 五、部署gopush-cluster
1.下载gopush-cluster及依赖包
```sh
$ go get -u github.com/Terry-Mao/gopush-cluster
$ go get -u github.com/Terry-Mao/goconf
$ go get -u github.com/garyburd/redigo/redis
$ go get -u code.google.com/p/go.net/websocket
$ go get -u launchpad.net/gozk/zookeeper
```
* 如果提示如下,说明需要安装谷歌的hg工具（安装mercurial,参考附资料1）

go: missing Mercurial command. See http://golang.org/s/gogetcmd
package code.google.com/p/go.net/websocket: exec: "hg": executable file not found in $PATH
* 如果提示如下,说明需要安装bzr工具（参考附资料2）

go: missing Bazaar command. See http://golang.org/s/gogetcmd
package launchpad.net/gozk/zookeeper: exec: "bzr": executable file not found in $PATH
* 如果提示如下,此时gozk已经下载下来了,需要修改gozk的cgo路径（参考附资料3）

launchpad.net/gozk/zookeeper
../zk.go:15:23: error: zookeeper.h: No such file or directory

3.安装message、comet、web模块
```sh
$ cd $GOPATH/src/github.com/Terry-Mao/gopush-cluster/message
$ go install
$ cp message
$ cp message.conf $GOPATH/bin/
$ cd ../comet/
$ go install
$ cp comet-example.conf /data/app/go/bin/
$ cd ../web/
$ go install
$ cp web.conf /data/app/go/bin/
```
到此所有的环境都搭建完成！
### 六、启动gopush-cluster
```sh
$ cd /$GOPATH/bin
$ nohup ./message -c message.conf &
$ nohup ./comet -c comet-example.conf &
$ nohup ./web -c web.conf &
```
* 如果报错如下(参考附资料4)

error while loading shared libraries: libzookeeper_mt.so.2: cannot open shared object file: No such file or directory
### 七、测试
1.推送公信（消息过期时间为expire=600秒）
```sh
$ curl -d "test2" http://localhost:8091/admin/push/public?expire=600
```
成功返回：{"msg":"ok","ret":0}
2.推送私信（消息过期时间为expire=600秒）
```sh
$ curl -d "test" http://localhost:8091/admin/push?key=Terry-Mao\&expire=600\&gid=0
```
成功返回：{“msg":"ok","ret":0}
3.获取离线消息接口
在浏览器中打开：http://localhost:8090/msg/get?key=Terry-Mao&mid=1&pmid=0
成功返回：
```json
{
    "data":{
        "msgs":[
            "{"msg":"test","expire":1391943609703654726,"mid":13919435497036558}"
        ],
        "pmsgs":[
            "{"msg":"test2","expire":1391943637016665915,"mid":13919435770166656}"
        ]
    },
    "msg":"ok",
    "ret":0
}
```
4.获取节点接口
在浏览器中打开：http://localhost:8090/server/get?key=Terry-Mao&proto=2
成功返回：
```json
{
    "data":{
        "server":"localhost:6969"
    },
    "msg":"ok",
    "ret":0
}
```
### 八、附资料
1.下载安装[hg](code.google.com/p/go.net/websocket)
```sh
$ wget http://mercurial.selenic.com/release/mercurial-1.4.1.tar.gz 
$ tar -xvf mercurial-1.4.1.tar.gz
$ cd mercurial-1.4.1
$ make
$ make install
```
2.下载安装bzr(只为下载包：launchpad.net/gozk/zookeeper)
```sh
$ yum install bzr.x86_64
```
3.修正cgo路径
```sh
$ vim $GOPATH/src/launchpad.net/gozk/zookeeper/zk.go
# 找到下面这行
# cgo CFLAGS: -I/usr/include/c-client-src -I/usr/include/zookeeper
# 更改为正确的路径:
# cgo CFLAGS: -I/usr/local/include/zookeeper
# 之后,不需要再尝试下载gozk包
```
4.报错：libzookeeper_mt.so.2无法找到
```sh
$ export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
$ sudo ldconfig
```
