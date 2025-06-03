// Copyright 2024 Multi-Range Request Support
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"bytes"
	"io/ioutil"
	"syscall"
	"testing"
)

// MockMultiRangeBackend implements StorageBackend for testing multi-range requests
type MockMultiRangeBackend struct {
	data []byte
	cap  Capabilities
}

func NewMockMultiRangeBackend(data []byte, supportsMultiRange bool) *MockMultiRangeBackend {
	return &MockMultiRangeBackend{
		data: data,
		cap: Capabilities{
			Name:               "mock",
			SupportsMultiRange: supportsMultiRange,
		},
	}
}

func (m *MockMultiRangeBackend) Init(key string) error                                                     { return nil }
func (m *MockMultiRangeBackend) Capabilities() *Capabilities                                               { return &m.cap }
func (m *MockMultiRangeBackend) Bucket() string                                                            { return "test-bucket" }
func (m *MockMultiRangeBackend) HeadBlob(param *HeadBlobInput) (*HeadBlobOutput, error)                   { return nil, nil }
func (m *MockMultiRangeBackend) ListBlobs(param *ListBlobsInput) (*ListBlobsOutput, error)                { return nil, nil }
func (m *MockMultiRangeBackend) DeleteBlob(param *DeleteBlobInput) (*DeleteBlobOutput, error)             { return nil, nil }
func (m *MockMultiRangeBackend) DeleteBlobs(param *DeleteBlobsInput) (*DeleteBlobsOutput, error)          { return nil, nil }
func (m *MockMultiRangeBackend) RenameBlob(param *RenameBlobInput) (*RenameBlobOutput, error)             { return nil, nil }
func (m *MockMultiRangeBackend) CopyBlob(param *CopyBlobInput) (*CopyBlobOutput, error)                   { return nil, nil }
func (m *MockMultiRangeBackend) PutBlob(param *PutBlobInput) (*PutBlobOutput, error)                      { return nil, nil }
func (m *MockMultiRangeBackend) MultipartBlobBegin(param *MultipartBlobBeginInput) (*MultipartBlobCommitInput, error) { return nil, nil }
func (m *MockMultiRangeBackend) MultipartBlobAdd(param *MultipartBlobAddInput) (*MultipartBlobAddOutput, error) { return nil, nil }
func (m *MockMultiRangeBackend) MultipartBlobAbort(param *MultipartBlobCommitInput) (*MultipartBlobAbortOutput, error) { return nil, nil }
func (m *MockMultiRangeBackend) MultipartBlobCommit(param *MultipartBlobCommitInput) (*MultipartBlobCommitOutput, error) { return nil, nil }
func (m *MockMultiRangeBackend) MultipartExpire(param *MultipartExpireInput) (*MultipartExpireOutput, error) { return nil, nil }
func (m *MockMultiRangeBackend) RemoveBucket(param *RemoveBucketInput) (*RemoveBucketOutput, error)       { return nil, nil }
func (m *MockMultiRangeBackend) MakeBucket(param *MakeBucketInput) (*MakeBucketOutput, error)             { return nil, nil }
func (m *MockMultiRangeBackend) Delegate() interface{}                                                     { return m }

func (m *MockMultiRangeBackend) GetBlob(param *GetBlobInput) (*GetBlobOutput, error) {
	start := param.Start
	end := start + param.Count
	if end > uint64(len(m.data)) {
		end = uint64(len(m.data))
	}
	
	data := m.data[start:end]
	return &GetBlobOutput{
		Body: ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}

func (m *MockMultiRangeBackend) GetBlobMultiRange(param *GetBlobMultiRangeInput) (*GetBlobMultiRangeOutput, error) {
	if !m.cap.SupportsMultiRange {
		return nil, syscall.ENOTSUP
	}

	var parts []MultiRangePart
	for _, r := range param.Ranges {
		start := r.Start
		end := start + r.Count
		if end > uint64(len(m.data)) {
			end = uint64(len(m.data))
		}
		
		data := m.data[start:end]
		parts = append(parts, MultiRangePart{
			Range:       r,
			ContentType: "application/octet-stream",
			Body:        ioutil.NopCloser(bytes.NewReader(data)),
		})
	}

	return &GetBlobMultiRangeOutput{
		Parts: parts,
	}, nil
}

func TestMultiRangeRequest(t *testing.T) {
	// Create test data
	testData := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	
	// Test with backend that supports multi-range
	backend := NewMockMultiRangeBackend(testData, true)
	
	ranges := []Range{
		{Start: 0, Count: 5},   // "01234"
		{Start: 10, Count: 5},  // "abcde"
		{Start: 20, Count: 5},  // "klmno"
	}
	
	resp, err := backend.GetBlobMultiRange(&GetBlobMultiRangeInput{
		Key:    "test-key",
		Ranges: ranges,
	})
	
	if err != nil {
		t.Fatalf("GetBlobMultiRange failed: %v", err)
	}
	
	if len(resp.Parts) != 3 {
		t.Fatalf("Expected 3 parts, got %d", len(resp.Parts))
	}
	
	// Verify each part
	expectedData := []string{"01234", "abcde", "klmno"}
	for i, part := range resp.Parts {
		data, err := ioutil.ReadAll(part.Body)
		if err != nil {
			t.Fatalf("Failed to read part %d: %v", i, err)
		}
		
		if string(data) != expectedData[i] {
			t.Errorf("Part %d: expected %q, got %q", i, expectedData[i], string(data))
		}
		
		part.Body.Close()
	}
}

func TestMultiRangeRequestUnsupported(t *testing.T) {
	// Test with backend that doesn't support multi-range
	backend := NewMockMultiRangeBackend([]byte("test"), false)
	
	ranges := []Range{
		{Start: 0, Count: 2},
	}
	
	_, err := backend.GetBlobMultiRange(&GetBlobMultiRangeInput{
		Key:    "test-key",
		Ranges: ranges,
	})
	
	if err != syscall.ENOTSUP {
		t.Fatalf("Expected ENOTSUP, got %v", err)
	}
}

func TestRangeAnalysis(t *testing.T) {
	// Test the gap ratio calculation logic
	ranges := []Range{
		{Start: 0, Count: 100},
		{Start: 200, Count: 100},   // 100 byte gap
		{Start: 400, Count: 100},   // 100 byte gap
	}
	
	totalRangeSize := uint64(0)
	for _, r := range ranges {
		totalRangeSize += r.Count
	}
	
	totalSpan := ranges[len(ranges)-1].Start + ranges[len(ranges)-1].Count - ranges[0].Start
	gapRatio := float64(totalSpan-totalRangeSize) / float64(totalSpan)
	
	expectedGapRatio := 200.0 / 500.0 // 200 bytes gap out of 500 total span
	if gapRatio != expectedGapRatio {
		t.Errorf("Expected gap ratio %f, got %f", expectedGapRatio, gapRatio)
	}
}

func TestS3BackendMultiRangeCapabilities(t *testing.T) {
	// Test that S3Backend reports multi-range support
	backend := &S3Backend{}
	backend.cap = Capabilities{
		Name:               "s3",
		MaxMultipartSize:   5 * 1024 * 1024 * 1024,
		SupportsMultiRange: true,
	}
	
	if !backend.Capabilities().SupportsMultiRange {
		t.Error("S3Backend should support multi-range requests")
	}
}

func TestGCSBackendMultiRangeCapabilities(t *testing.T) {
	// Test that GCS Backend reports multi-range support
	backend := &GCSBackend{}
	backend.cap = Capabilities{
		MaxMultipartSize:   5 * 1024 * 1024 * 1024,
		Name:               "gcs",
		NoParallelMultipart: true,
		SupportsMultiRange: true,
	}
	
	if !backend.Capabilities().SupportsMultiRange {
		t.Error("GCSBackend should support multi-range requests")
	}
}

func TestAzureBlobBackendMultiRangeCapabilities(t *testing.T) {
	// Test that Azure Blob Backend reports no multi-range support
	backend := &AZBlob{}
	backend.cap = Capabilities{
		MaxMultipartSize:   100 * 1024 * 1024,
		Name:               "wasb",
		SupportsMultiRange: false,
	}
	
	if backend.Capabilities().SupportsMultiRange {
		t.Error("AZBlob should not support multi-range requests")
	}
}
