# qmv
简单的七牛空间文件批量移动工具

## 使用场景
自己的 bucket 没绑定域名，里面存了很多文件，由于七牛回收了测试域名而导致文件无法下载。此时，可以新建一个 bucket，然后通过这个工具把被回收了域名的 bucket 中的文件全部移动到新 bucket 中，这样就能下载了。

## 编译
`go build -o qmv main.go`

## 使用方法
设置好环境变量 `Q_AK` 和 `Q_SK`（登录 [七牛 Portal](portal.qiniu.com)，点击头像 -> 密钥管理），然后运行 `./qmv <srcBucket> <destBucket>` 即可。
