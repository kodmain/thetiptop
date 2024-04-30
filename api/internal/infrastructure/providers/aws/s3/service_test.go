package s3_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/aws"
	service "github.com/kodmain/thetiptop/api/internal/infrastructure/providers/aws/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3API struct {
	mock.Mock
}

func (m *MockS3API) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func TestNew(t *testing.T) {
	svc, err := service.New()
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	svc, err = service.New()
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestGetObject(t *testing.T) {
	mockS3 := new(MockS3API)
	service := &service.Service{
		API: mockS3,
	}

	bucket := "test-bucket"
	item := "test-item"
	content := "Hello, world!"

	// Set up the expected response
	mockS3.On("GetObject", aws.CTX, mock.AnythingOfType("*s3.GetObjectInput"), mock.Anything).Return(&s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte(content))),
	}, nil)

	// Call the function
	result, err := service.GetObject(&bucket, &item)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, content, result.String())
	mockS3.AssertExpectations(t)
}

/*
func TestGetObject(t *testing.T) {
	bucket := "test-bucket"
	item := "test-item"

	expectedData := []byte("test data")
	expectedBuffer := bytes.NewBuffer(expectedData)

	// Mock the s3.Client and its GetObject method

	mockClient := &s3.MockClient{
		GetObjectFunc: func(input *s3.GetObjectInput) s3.GetObjectRequest {
			if *input.Bucket != bucket || *input.Key != item {
				t.Errorf("Unexpected input: bucket=%s, item=%s", *input.Bucket, *input.Key)
			}

			return s3.GetObjectRequest{
				Request: &aws.Request{
					Data: expectedData,
				},
			}
		},
	}

	// Call the GetObject function
	buffer, err := s3.GetObject(&bucket, &item)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the returned buffer matches the expected buffer
	if !bytes.Equal(buffer.Bytes(), expectedBuffer.Bytes()) {
		t.Errorf("Unexpected buffer: got=%v, want=%v", buffer.Bytes(), expectedBuffer.Bytes())
	}
}
*/
