<img src="doc/goofys.png" height="32" width="32" align="middle" /> Goofys 是一个高性能、类POSIX的 [Amazon S3](https://aws.amazon.com/s3/) 文件系统，使用Go语言编写

[![Build Status](https://travis-ci.org/kahing/goofys.svg?branch=master)](https://travis-ci.org/kahing/goofys)
[![Github All Releases](https://img.shields.io/github/downloads/kahing/goofys/total.svg)](https://github.com/kahing/goofys/releases/)
[![Twitter Follow](https://img.shields.io/twitter/follow/s3goofys.svg?style=social&label=Follow)](https://twitter.com/s3goofys)
[![Stack Overflow Questions](https://img.shields.io/stackexchange/stackoverflow/t/goofys?label=Stack%20Overflow%20questions)](https://stackoverflow.com/search?q=%5Bgoofys%5D+is%3Aquestion)

# 概述

Goofys 允许您将S3存储桶挂载为文件系统。

它是一个"文件式系统"而不是"文件系统"，因为goofys优先考虑性能，其次才是POSIX兼容性。特别是那些在S3上难以支持或需要多次往返的操作要么会失败（随机写入），要么被模拟（没有单文件权限）。Goofys没有磁盘数据缓存（请查看[catfs](https://github.com/kahing/catfs)），一致性模型是close-to-open。

## 性能特性

Goofys包含多项性能优化：

* **多范围请求**：启用后，goofys可以在单个请求中发出HTTP多范围请求来获取多个非连续字节范围，减少稀疏文件访问模式的延迟。
* **预读缓冲**：智能预取数据以提高顺序读取性能。
* **并行操作**：并发处理多个文件操作。

# 安装

* 在Linux上，通过[预构建二进制文件](https://github.com/kahing/goofys/releases/latest/download/goofys)安装。
如果您想在启动时挂载，可能还需要安装fuse。

* 在macOS上，通过[Homebrew](https://brew.sh/)安装：

```ShellSession
$ brew cask install osxfuse
$ brew install goofys
```

* 或使用Go 1.10或更高版本从源码构建：

```ShellSession
$ export GOPATH=$HOME/work
$ go get github.com/kahing/goofys
$ go install github.com/kahing/goofys
```

# 使用方法

```ShellSession
$ cat ~/.aws/credentials
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
$ $GOPATH/bin/goofys <bucket> <mountpoint>
$ $GOPATH/bin/goofys <bucket:prefix> <mountpoint> # 如果您只想挂载前缀下的对象
```

用户也可以通过[AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html)或`AWS_ACCESS_KEY_ID`和`AWS_SECRET_ACCESS_KEY`环境变量配置凭据。

要在启动时挂载S3存储桶，请确保为`root`配置了凭据，并可以将此添加到`/etc/fstab`：

```
goofys#bucket   /mnt/mountpoint        fuse     _netdev,allow_other,--file-mode=0666,--dir-mode=0777    0       0
```

## 性能调优

对于稀疏文件访问模式的性能改进，您可以启用多范围请求：

```ShellSession
$ $GOPATH/bin/goofys --enable-multi-range <bucket> <mountpoint>
```

其他多范围选项：
* `--multi-range-batch-size N`：每个请求的最大范围数（默认：5）
* `--multi-range-threshold N`：触发多范围的最小间隙大小（字节）（默认：1048576）

**注意**：多范围请求在AWS S3和Google Cloud Storage上受支持，但在Azure存储服务上不受支持。

另请参阅：[Azure Blob Storage、Azure Data Lake Gen1和Azure Data Lake Gen2的说明](https://github.com/kahing/goofys/blob/master/README-azure.md)。

有更多问题？查看[其他人提出的问题](https://github.com/kahing/goofys/issues?utf8=%E2%9C%93&q=is%3Aissue%20label%3Aquestion%20)

# 基准测试

以下是在`c4.xlarge`实例上运行的一些基准测试结果，该实例连接到同一区域的存储桶。单位为秒。对于[s3fs](https://github.com/s3fs-fuse/s3fs-fuse)，缓存目录被清除以模拟冷读取。

![基准测试结果](/bench/bench.png?raw=true "基准测试")

要运行基准测试，配置EC2实例角色以能够写入`$TESTBUCKET`，然后执行：
```ShellSession
$ sudo docker run -e BUCKET=$TESTBUCKET -e CACHE=false --rm --privileged --net=host -v /tmp/cache:/tmp/cache kahing/goofys-bench
# 结果将写入$TESTBUCKET
```

另请参阅：[缓存基准测试结果](https://github.com/kahing/goofys/blob/master/bench/cache/README.md)和[Azure上的结果](https://github.com/kahing/goofys/blob/master/bench/azure/README.md)。

# 许可证

版权所有 (C) 2015 - 2019 Ka-Hing Cheung

根据Apache License 2.0许可

# 当前状态

goofys已在Linux和macOS下测试。

## 限制

* 只能追加写入
* 不支持随机写入
* 不支持符号链接
* `mtime`反映服务器端时间戳
* 不能`rename`非空目录
* 不支持硬链接
* `fsync`被忽略，文件只在`close`时刷新

## 与非AWS S3的兼容性

goofys已与以下非AWS S3提供商测试：

* Amplidata / WD ActiveScale
* Ceph（例如：Digital Ocean Spaces、DreamObjects、gridscale）
* EdgeFS
* EMC Atmos
* Google Cloud Storage（支持多范围请求）
* Minio（有限）
* OpenStack Swift
* S3Proxy
* Scaleway
* Wasabi

此外，goofys还可以与以下非S3对象存储一起使用：

* Azure Blob Storage（不支持多范围请求）
* Azure Data Lake Gen1（不支持多范围请求）
* Azure Data Lake Gen2（不支持多范围请求）

# 参考资料

  * 数据存储在[Amazon S3](https://aws.amazon.com/s3/)
  * [Amazon SDK for Go](https://github.com/aws/aws-sdk-go)
  * 其他相关的fuse文件系统
    * [catfs](https://github.com/kahing/catfs)：可与goofys一起使用的缓存层
    * [s3fs](https://github.com/s3fs-fuse/s3fs-fuse)：另一个流行的S3文件系统
    * [gcsfuse](https://github.com/googlecloudplatform/gcsfuse)：
      [Google Cloud Storage](https://cloud.google.com/storage/)的文件系统。Goofys
      从这个项目借用了一些骨架代码。
  * [S3Proxy](https://github.com/andrewgaul/s3proxy)用于`go test`
  * [fuse绑定](https://github.com/jacobsa/fuse)，也被`gcsfuse`使用
