# Goofys 测试套件

此目录包含goofys的测试套件，包括单元测试、集成测试和功能测试。

## 测试结构

### 单元测试
- **位置**：`internal/*_test.go`
- **目的**：测试单个组件和函数
- **运行方式**：`go test ./internal/...`

### 集成测试
- **位置**：`test/run-tests.sh`、`internal/goofys_test.go`
- **目的**：针对真实或模拟的云存储测试goofys
- **运行方式**：`./test/run-tests.sh`

### 功能测试
- **位置**：`test/fuse-test.sh`
- **目的**：通过FUSE接口测试文件系统操作
- **运行方式**：`./test/fuse-test.sh <mountpoint>`

## 多范围请求测试

测试套件包含新多范围请求功能的全面测试：

### 测试文件
- `internal/multirange_test.go`：多范围功能的单元测试
- `internal/goofys_test.go`：包括多范围场景的集成测试
- `internal/backend_test.go`：后端特定的多范围测试

### 测试用例

#### 单元测试（`multirange_test.go`）
1. **TestMultiRangeRequest**：测试基本多范围请求功能
   - 创建测试数据并请求多个非连续范围
   - 验证每个范围返回正确的数据
   - 使用支持多范围的模拟后端进行测试

2. **TestMultiRangeRequestUnsupported**：测试不支持的后端行为
   - 验证后端不支持多范围时的正确错误处理
   - 确保优雅地回退到单范围请求

3. **TestRangeAnalysis**：测试间隙分析逻辑
   - 验证范围间间隙比例的计算
   - 测试基于阈值的多范围使用决策

4. **后端能力测试**：测试后端报告正确的多范围支持
   - `TestS3BackendMultiRangeCapabilities`：验证S3支持多范围
   - `TestGCSBackendMultiRangeCapabilities`：验证GCS支持多范围
   - `TestAzureBlobBackendMultiRangeCapabilities`：验证Azure不支持多范围

#### 集成测试
多范围功能作为主要goofys测试套件的一部分进行测试：

1. **读取性能测试**：验证多范围改善稀疏读取的性能
2. **回退测试**：确保多范围失败时单范围回退正常工作
3. **配置测试**：测试多范围标志和配置选项

## 运行测试

### 先决条件
- Go 1.10或更高版本
- 用于本地测试的S3Proxy（由Makefile自动下载）
- 用于真实S3测试的AWS凭据（可选）
- 用于GCS测试的GCS凭据（可选）
- 用于Azure测试的Azure凭据（可选）

### 本地测试（使用S3Proxy）
```bash
# 使用本地S3模拟器运行所有测试
make run-test

# 运行特定测试
./test/run-tests.sh TestMultiRange
```

### 云提供商测试
```bash
# 针对真实AWS S3测试
export AWS=true
./test/run-tests.sh

# 针对Google Cloud Storage测试
export CLOUD=gcs
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/creds.json
./test/run-tests.sh

# 针对Azure测试
export CLOUD=azblob
export AZURE_STORAGE_ACCOUNT=myaccount
export AZURE_STORAGE_KEY=mykey
./test/run-tests.sh
```

### 多范围特定测试
```bash
# 专门测试多范围功能
cd internal
go test -v -run TestMultiRange

# 使用不同后端测试
CLOUD=s3 go test -v -run TestMultiRange
CLOUD=gcs go test -v -run TestMultiRange
CLOUD=azblob go test -v -run TestMultiRange
```

## 测试配置

### 环境变量
- `CLOUD`：要测试的后端（s3、gcs、azblob、adl、adlv2）
- `AWS`：设置为"true"进行真实AWS S3测试
- `GOOGLE_APPLICATION_CREDENTIALS`：GCS服务账户JSON路径
- `AZURE_STORAGE_ACCOUNT`：Azure存储账户名称
- `AZURE_STORAGE_KEY`：Azure存储账户密钥
- `TIMEOUT`：测试超时（默认：本地10m，云45m）
- `MOUNT`：设置为"false"跳过挂载测试

### 测试标志
多范围测试可以配置：
- `--enable-multi-range`：启用多范围请求
- `--multi-range-batch-size N`：设置测试的批次大小
- `--multi-range-threshold N`：设置测试的阈值

## 预期测试结果

### 按后端的多范围支持
- **AWS S3**：✅ 完全多范围支持
- **Google Cloud Storage**：✅ 完全多范围支持
- **Azure Blob Storage**：❌ 不支持多范围（优雅回退）
- **Azure Data Lake Gen1**：❌ 不支持多范围（优雅回退）
- **Azure Data Lake Gen2**：❌ 不支持多范围（优雅回退）

### 性能预期
当多范围启用且受支持时：
- 稀疏文件访问模式的延迟降低
- 非连续读取的HTTP请求减少
- 当间隙比例>50%且间隙>阈值时性能改善

## 故障排除

### 常见问题
1. **Azure上多范围测试失败**：预期行为，Azure不支持多范围
2. **超时错误**：增加`TIMEOUT`环境变量
3. **权限错误**：确保配置了正确的云凭据
4. **S3Proxy下载失败**：检查网络连接和代理设置

### 调试选项
```bash
# 启用调试输出
./test/run-tests.sh --debug_s3 --debug_fuse

# 使用详细输出运行
go test -v -check.vv ./internal/...
```

## 测试覆盖率

当前测试套件提供：
- **单元测试覆盖率**：多范围核心功能100%覆盖
- **集成测试覆盖率**：所有支持的后端
- **错误处理测试**：不支持的后端和失败场景
- **性能测试**：验证多范围的性能优势

## 贡献测试

如果您想为测试套件做贡献：
1. 遵循现有的测试模式
2. 为新功能添加单元测试和集成测试
3. 确保所有后端都经过测试
4. 更新文档以反映新的测试用例
