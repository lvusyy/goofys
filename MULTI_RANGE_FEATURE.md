# Multi-Range Request Feature Documentation

## Overview

This document describes the implementation and documentation updates for the HTTP multi-range request feature in goofys. This feature allows goofys to fetch multiple non-contiguous byte ranges in a single HTTP request, significantly improving performance for sparse file access patterns.

## Feature Implementation

### Core Components

1. **Backend Interface Extensions** (`internal/backend.go`)
   - Added `SupportsMultiRange` capability flag
   - Added `Range` struct for defining byte ranges
   - Added `GetBlobMultiRangeInput` and `GetBlobMultiRangeOutput` types
   - Added `GetBlobMultiRange` method to `StorageBackend` interface

2. **S3 Backend Support** (`internal/backend_s3.go`)
   - Implemented `GetBlobMultiRange` method
   - Added multipart response parsing
   - Set `SupportsMultiRange: true` in capabilities

3. **GCS Backend Support** (`internal/backend_gcs.go`)
   - Implemented `GetBlobMultiRange` method using individual range readers
   - Set `SupportsMultiRange: true` in capabilities

4. **Azure Backend Limitations**
   - Azure Blob Storage (`internal/backend_azblob.go`): `SupportsMultiRange: false`
   - Azure Data Lake Gen1 (`internal/backend_adlv1.go`): `SupportsMultiRange: false`
   - Azure Data Lake Gen2 (`internal/backend_adlv2.go`): `SupportsMultiRange: false`

5. **File Handling** (`internal/file.go`)
   - Added `S3MultiRangeReadBuffer` for managing multi-range responses
   - Added `analyzeReadPattern` for detecting beneficial multi-range scenarios
   - Added `readAheadMultiRange` for executing multi-range read-ahead
   - Integrated multi-range logic into existing read-ahead system

6. **Configuration** (`api/common/config.go`, `internal/flags.go`)
   - Added `EnableMultiRange` flag (default: false)
   - Added `MultiRangeBatchSize` flag (default: 5)
   - Added `MultiRangeThreshold` flag (default: 1MB)

### Algorithm Details

The multi-range feature uses intelligent analysis to determine when to use multi-range requests:

1. **Gap Analysis**: Calculates the ratio of gaps to total span across ranges
2. **Threshold Check**: Only triggers when gaps exceed the configured threshold
3. **Batch Limiting**: Limits the number of ranges per request to avoid oversized requests
4. **Fallback**: Gracefully falls back to single-range requests on failure

## Documentation Updates

### 1. Main README (`README.md`)

**Added:**
- Performance Features section describing multi-range requests
- Performance Tuning section with usage examples
- Multi-range support notes in compatibility sections

**Key Changes:**
- Added multi-range description in Overview
- Added usage examples with `--enable-multi-range` flag
- Updated compatibility lists to indicate multi-range support status

### 2. Azure Documentation (`README-azure.md`)

**Added:**
- Multi-range limitation notes for all Azure services
- Clear statements that multi-range flags are ignored on Azure

### 3. GCS Documentation (`README-gcs.md`)

**Added:**
- Multi-range support announcement
- Usage example with multi-range enabled

### 4. Test Documentation (`test/README.md`)

**Created comprehensive test documentation including:**
- Test structure overview
- Multi-range specific test cases
- Running instructions for different scenarios
- Expected results by backend
- Troubleshooting guide

### 5. Build System (`Makefile`)

**Added:**
- `run-test-multirange` target for running multi-range specific tests

## Test Coverage

### Unit Tests (`internal/multirange_test.go`)

1. **TestMultiRangeRequest**: Basic functionality test
2. **TestMultiRangeRequestUnsupported**: Unsupported backend handling
3. **TestRangeAnalysis**: Gap analysis algorithm testing
4. **Backend Capability Tests**: Verify correct capability reporting

### Integration Tests

Multi-range functionality is tested as part of the main goofys test suite with various backends.

## Usage Examples

### Basic Usage
```bash
# Enable multi-range requests
goofys --enable-multi-range mybucket /mnt/mybucket
```

### Advanced Configuration
```bash
# Custom batch size and threshold
goofys --enable-multi-range \
       --multi-range-batch-size 10 \
       --multi-range-threshold 2097152 \
       mybucket /mnt/mybucket
```

### Backend-Specific Usage

**AWS S3:**
```bash
goofys --enable-multi-range s3://mybucket /mnt/mybucket
```

**Google Cloud Storage:**
```bash
GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" \
goofys --enable-multi-range gs://mybucket /mnt/mybucket
```

**Azure (multi-range ignored):**
```bash
goofys --enable-multi-range wasb://container /mnt/container
# Note: Multi-range flag is ignored, falls back to single ranges
```

## Performance Benefits

### When Multi-Range Helps
- Sparse file access patterns
- Reading multiple small sections of large files
- Applications that seek frequently within files
- Workloads with significant gaps between read ranges

### Performance Metrics
- Reduced HTTP request count for non-contiguous reads
- Lower latency for sparse access patterns
- Improved throughput when gap ratio > 50%

### When NOT to Use Multi-Range
- Sequential file reading (no benefit)
- Small files (overhead not justified)
- Backends that don't support it (Azure)

## Backend Support Matrix

| Backend | Multi-Range Support | Notes |
|---------|-------------------|-------|
| AWS S3 | ✅ Yes | Full HTTP multi-range support |
| Google Cloud Storage | ✅ Yes | Implemented via multiple range readers |
| Azure Blob Storage | ❌ No | Not supported by Azure API |
| Azure Data Lake Gen1 | ❌ No | Not supported by Azure API |
| Azure Data Lake Gen2 | ❌ No | Not supported by Azure API |
| Minio | ⚠️ Unknown | Depends on S3 compatibility |
| Other S3-compatible | ⚠️ Varies | Depends on implementation |

## Configuration Reference

### Command Line Flags

- `--enable-multi-range`: Enable multi-range requests (default: false)
- `--multi-range-batch-size N`: Max ranges per request (default: 5)
- `--multi-range-threshold N`: Min gap size in bytes (default: 1048576)

### Environment Variables

Multi-range settings are configured via command-line flags only.

### Tuning Guidelines

- **Batch Size**: Start with 5, increase for very sparse access patterns
- **Threshold**: Start with 1MB, adjust based on typical gap sizes
- **Enable**: Only enable if you have sparse access patterns

## Troubleshooting

### Common Issues

1. **Multi-range not working on Azure**: Expected, Azure doesn't support it
2. **No performance improvement**: Check if access pattern is actually sparse
3. **Increased latency**: May occur with small gaps, increase threshold

### Debug Options

```bash
# Enable debug output to see multi-range decisions
goofys --debug_s3 --enable-multi-range mybucket /mnt/mybucket
```

## Future Enhancements

Potential improvements for future versions:
- Adaptive threshold adjustment based on access patterns
- Better integration with read-ahead algorithms
- Support for more S3-compatible backends
- Performance metrics and monitoring
