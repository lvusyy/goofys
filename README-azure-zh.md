# Azure Blob Storage

```ShellSession
$ cat ~/.azure/config
[storage]
account = "myblobstorage"
key = "MY-STORAGE-KEY"
$ $GOPATH/bin/goofys wasb://container <mountpoint>
$ $GOPATH/bin/goofys wasb://container:prefix <mountpoint> # 如果您只想挂载前缀下的对象
```

**注意**：Azure Blob Storage不支持HTTP多范围请求。使用Azure后端时，`--enable-multi-range`标志将被忽略。

用户也可以通过`AZURE_STORAGE_ACCOUNT`和`AZURE_STORAGE_KEY`环境变量配置凭据。详情请参阅[Azure CLI配置](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest#cli-configuration-values-and-environment-variables)。Goofys尚不支持`connection_string`或`sas_token`。

Goofys也接受完整的`wasb` URI：
```ShellSession
$ $GOPATH/bin/goofys wasb://container@myaccount.blob.core.windows.net <mountpoint>
$ $GOPATH/bin/goofys wasb://container@myaccount.blob.core.windows.net/prefix <mountpoint>
```

在这种情况下，可以省略`~/.azure/config`中的账户配置或`AZURE_STORAGE_ACCOUNT`。或者，也可以使用`--endpoint`来指定存储账户：

```ShellSession
$ $GOPATH/bin/goofys --endpoint https://myaccount.blob.core.windows.net wasb://container <mountpoint>
$ $GOPATH/bin/goofys --endpoint https://myaccount.blob.core.windows.net wasb://container:prefix <mountpoint>
```

注意，如果未指定完整的`wasb` URI，前缀分隔符是`:`。

最后，除了指定存储账户访问密钥外，goofys还可以使用[Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)访问令牌：

```ShellSession
$ az login
# 列出所有订阅并选择需要的订阅
$ az account list
# 选择当前订阅（从上一步获取其ID）
$ az account set --subscription <name or id>
$ $GOPATH/bin/goofys wasb://container@myaccount.blob.core.windows.net <mountpoint>
```

# Azure Data Lake Storage Gen1

按照上述Azure CLI登录序列，然后：

```ShellSession
$ $GOPATH/bin/goofys adl://servicename.azuredatalakestore.net <mountpoint>
$ $GOPATH/bin/goofys adl://servicename.azuredatalakestore.net:prefix <mountpoint>
```

**注意**：Azure Data Lake Gen1不支持HTTP多范围请求。

# Azure Data Lake Storage Gen2

按照上述[Azure Blob Storage](https://github.com/kahing/goofys/blob/master/README-azure.md#azure-blob-storage)的方式配置凭据，然后：

```ShellSession
$ $GOPATH/bin/goofys abfs://container <mountpoint>
$ $GOPATH/bin/goofys abfs://container:prefix <mountpoint>
```

**注意**：Azure Data Lake Gen2不支持HTTP多范围请求。

## Azure服务的多范围请求限制

所有Azure存储服务（Blob Storage、Data Lake Gen1、Data Lake Gen2）都不支持HTTP多范围请求功能。这意味着：

- `--enable-multi-range`标志在使用Azure后端时将被忽略
- Goofys将自动回退到单范围请求
- 对于稀疏文件访问模式，性能可能不如支持多范围的后端（如AWS S3或Google Cloud Storage）

如果您的工作负载需要多范围请求的性能优势，建议考虑使用AWS S3或Google Cloud Storage。
