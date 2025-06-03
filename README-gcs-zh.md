# Google Cloud Storage (GCS)

Google Cloud Storage支持HTTP多范围请求，可以改善稀疏文件访问模式的性能。使用`--enable-multi-range`标志启用此功能。

## 先决条件

服务账户凭据或用户身份验证。确保服务账户或用户对GCS下的存储桶/对象具有适当的权限。

要成功挂载，我们要求用户对存储桶具有对象列表（`storage.objects.list`）权限。

### 服务账户凭据

创建服务账户凭据（https://cloud.google.com/iam/docs/creating-managing-service-accounts）并生成JSON凭据文件。

### 用户身份验证和`gcloud`默认身份验证
用户可以通过首先安装cloud sdk（https://cloud.google.com/sdk/）并运行`gcloud auth application-default login`命令来向gcloud的默认环境进行身份验证。

## 使用Goofys for GCS

### 使用服务账户凭据文件
```
GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" goofys gs://[BUCKET] /path/to/mount
```

### 使用用户身份验证（`gcloud auth application-default login`）

```
goofys gs://[BUCKET] [MOUNT DIRECTORY]
```

### 启用多范围请求

```
GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" goofys --enable-multi-range gs://[BUCKET] /path/to/mount
```

## 多范围请求的优势

在Google Cloud Storage上启用多范围请求可以提供以下性能优势：

### 适用场景
- **稀疏文件访问**：当应用程序需要读取文件的多个非连续部分时
- **随机访问模式**：频繁在文件内跳转的应用程序
- **大文件的部分读取**：只需要读取大文件中的特定段落

### 性能改进
- **减少HTTP请求数量**：将多个范围请求合并为单个请求
- **降低延迟**：减少网络往返次数
- **提高吞吐量**：当间隙比例大于50%时效果显著

### 配置选项
```ShellSession
# 基本启用
goofys --enable-multi-range gs://mybucket /mnt/mybucket

# 自定义批次大小和阈值
goofys --enable-multi-range \
       --multi-range-batch-size 10 \
       --multi-range-threshold 2097152 \
       gs://mybucket /mnt/mybucket
```

### 配置参数说明
- `--multi-range-batch-size`：单个多范围请求中包含的最大范围数（默认：5）
- `--multi-range-threshold`：触发多范围请求的最小间隙大小（字节）（默认：1048576）

## 最佳实践

1. **评估访问模式**：多范围请求最适合稀疏访问模式
2. **调整阈值**：根据您的典型间隙大小调整阈值
3. **监控性能**：使用`--debug_s3`标志查看多范围决策
4. **测试配置**：在生产环境中使用前测试不同的批次大小

## 故障排除

如果多范围请求没有提供预期的性能改进：
1. 检查您的访问模式是否真的是稀疏的
2. 尝试调整`--multi-range-threshold`值
3. 使用调试标志查看是否正在使用多范围请求
4. 考虑您的网络延迟和带宽特性
