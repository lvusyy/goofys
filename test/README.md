# Goofys Test Suite

This directory contains the test suite for goofys, including unit tests, integration tests, and functional tests.

## Test Structure

### Unit Tests
- **Location**: `internal/*_test.go`
- **Purpose**: Test individual components and functions
- **Run with**: `go test ./internal/...`

### Integration Tests
- **Location**: `test/run-tests.sh`, `internal/goofys_test.go`
- **Purpose**: Test goofys against real or emulated cloud storage
- **Run with**: `./test/run-tests.sh`

### Functional Tests
- **Location**: `test/fuse-test.sh`
- **Purpose**: Test filesystem operations through FUSE interface
- **Run with**: `./test/fuse-test.sh <mountpoint>`

## Multi-Range Request Tests

The test suite includes comprehensive tests for the new multi-range request feature:

### Test Files
- `internal/multirange_test.go`: Unit tests for multi-range functionality
- `internal/goofys_test.go`: Integration tests including multi-range scenarios
- `internal/backend_test.go`: Backend-specific multi-range tests

### Test Cases

#### Unit Tests (`multirange_test.go`)
1. **TestMultiRangeRequest**: Tests basic multi-range request functionality
   - Creates test data and requests multiple non-contiguous ranges
   - Verifies correct data is returned for each range
   - Tests with mock backend that supports multi-range

2. **TestMultiRangeRequestUnsupported**: Tests behavior with unsupported backends
   - Verifies proper error handling when backend doesn't support multi-range
   - Ensures graceful fallback to single-range requests

3. **TestRangeAnalysis**: Tests the gap analysis logic
   - Verifies calculation of gap ratios between ranges
   - Tests threshold-based decision making for multi-range usage

4. **Backend Capability Tests**: Tests that backends report correct multi-range support
   - `TestS3BackendMultiRangeCapabilities`: Verifies S3 supports multi-range
   - `TestGCSBackendMultiRangeCapabilities`: Verifies GCS supports multi-range
   - `TestAzureBlobBackendMultiRangeCapabilities`: Verifies Azure does not support multi-range

#### Integration Tests
Multi-range functionality is tested as part of the main goofys test suite:

1. **Read Performance Tests**: Verify multi-range improves performance for sparse reads
2. **Fallback Tests**: Ensure single-range fallback works when multi-range fails
3. **Configuration Tests**: Test multi-range flags and configuration options

## Running Tests

### Prerequisites
- Go 1.10 or later
- S3Proxy for local testing (automatically downloaded by Makefile)
- AWS credentials for real S3 testing (optional)
- GCS credentials for GCS testing (optional)
- Azure credentials for Azure testing (optional)

### Local Testing (with S3Proxy)
```bash
# Run all tests with local S3 emulator
make run-test

# Run specific test
./test/run-tests.sh TestMultiRange
```

### Cloud Provider Testing
```bash
# Test against real AWS S3
export AWS=true
./test/run-tests.sh

# Test against Google Cloud Storage
export CLOUD=gcs
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/creds.json
./test/run-tests.sh

# Test against Azure
export CLOUD=azblob
export AZURE_STORAGE_ACCOUNT=myaccount
export AZURE_STORAGE_KEY=mykey
./test/run-tests.sh
```

### Multi-Range Specific Testing
```bash
# Test multi-range functionality specifically
cd internal
go test -v -run TestMultiRange

# Test with different backends
CLOUD=s3 go test -v -run TestMultiRange
CLOUD=gcs go test -v -run TestMultiRange
CLOUD=azblob go test -v -run TestMultiRange
```

## Test Configuration

### Environment Variables
- `CLOUD`: Backend to test (s3, gcs, azblob, adl, adlv2)
- `AWS`: Set to "true" for real AWS S3 testing
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to GCS service account JSON
- `AZURE_STORAGE_ACCOUNT`: Azure storage account name
- `AZURE_STORAGE_KEY`: Azure storage account key
- `TIMEOUT`: Test timeout (default: 10m for local, 45m for cloud)
- `MOUNT`: Set to "false" to skip mount tests

### Test Flags
Multi-range tests can be configured with:
- `--enable-multi-range`: Enable multi-range requests
- `--multi-range-batch-size N`: Set batch size for testing
- `--multi-range-threshold N`: Set threshold for testing

## Expected Test Results

### Multi-Range Support by Backend
- **AWS S3**: ✅ Full multi-range support
- **Google Cloud Storage**: ✅ Full multi-range support  
- **Azure Blob Storage**: ❌ Multi-range not supported (graceful fallback)
- **Azure Data Lake Gen1**: ❌ Multi-range not supported (graceful fallback)
- **Azure Data Lake Gen2**: ❌ Multi-range not supported (graceful fallback)

### Performance Expectations
When multi-range is enabled and supported:
- Reduced latency for sparse file access patterns
- Fewer HTTP requests for non-contiguous reads
- Improved performance when gap ratio > 50% and gaps > threshold

## Troubleshooting

### Common Issues
1. **Multi-range tests fail on Azure**: Expected behavior, Azure doesn't support multi-range
2. **Timeout errors**: Increase `TIMEOUT` environment variable
3. **Permission errors**: Ensure proper cloud credentials are configured
4. **S3Proxy download fails**: Check internet connection and proxy settings

### Debug Options
```bash
# Enable debug output
./test/run-tests.sh --debug_s3 --debug_fuse

# Run with verbose output
go test -v -check.vv ./internal/...
```
