package main

import (
	"context"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/napmany/nakamasrv/mockobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"testing/fstest"
)

func TestMain(m *testing.M) {
	filestorage = fstest.MapFS{
		"core/1.0.0.json": {Data: []byte(`{
  "some": "data",
  "more": "data",
  "even": "more data",
  "and": "even more data"
}`)},
	}

	os.Exit(m.Run())
}

func TestRpcTransferFile(t *testing.T) {
	validResponse := `{"type":"core","version":"1.0.0","hash":"cbfab3df1f0156ba9eb8e292b754b8cd4f802582ce44b0a0551e918cf3d09092","content":{"some":"data","more":"data","even":"more data","and":"even more data"}}`
	tests := []struct {
		name             string
		payload          string
		expectedError    error
		expectedResponse string
		expectedCalls    func(*mockobject.NakamaModuleMock)
	}{
		{
			"valid request with matching hash, store item in db",
			`{"type":"core","version":"1.0.0"}`,
			nil,
			validResponse,
			func(nkMock *mockobject.NakamaModuleMock) {
				nkMock.On("StorageRead", context.Background(), []*runtime.StorageRead{{
					Collection: filestorageCollection,
					Key:        "core-1.0.0",
				}}).
					Return([]*api.StorageObject{}, nil)
				nkMock.On("StorageWrite", context.Background(), []*runtime.StorageWrite{
					{
						Collection:      filestorageCollection,
						Key:             "core-1.0.0",
						Value:           validResponse,
						PermissionRead:  runtime.STORAGE_PERMISSION_OWNER_READ,
						PermissionWrite: runtime.STORAGE_PERMISSION_OWNER_WRITE,
					},
				}).
					Return([]*api.StorageObjectAck{}, nil)
			},
		},
		{
			"empty payload, use defaults",
			`{}`,
			nil,
			validResponse,
			func(nkMock *mockobject.NakamaModuleMock) {
				nkMock.On("StorageRead", context.Background(), []*runtime.StorageRead{{
					Collection: filestorageCollection,
					Key:        "core-1.0.0",
				}}).
					Return([]*api.StorageObject{}, nil)
				nkMock.On("StorageWrite", context.Background(), []*runtime.StorageWrite{
					{
						Collection:      filestorageCollection,
						Key:             "core-1.0.0",
						Value:           validResponse,
						PermissionRead:  runtime.STORAGE_PERMISSION_OWNER_READ,
						PermissionWrite: runtime.STORAGE_PERMISSION_OWNER_WRITE,
					},
				}).
					Return([]*api.StorageObjectAck{}, nil)
			},
		},
		{
			"valid request with matching hash, return item from db",
			`{"type":"core","version":"1.0.0"}`,
			nil,
			validResponse,
			func(nkMock *mockobject.NakamaModuleMock) {
				nkMock.On("StorageRead", context.Background(), mock.Anything).
					Return([]*api.StorageObject{{Value: validResponse}}, nil)
			},
		},
		{
			"non-matching hash",
			`{"type":"core","version":"1.0.0","hash":"invalidhash"}`,
			nil,
			`{"type":"core","version":"1.0.0","hash":"cbfab3df1f0156ba9eb8e292b754b8cd4f802582ce44b0a0551e918cf3d09092","content":null}`,
			func(nkMock *mockobject.NakamaModuleMock) {
				nkMock.On("StorageRead", context.Background(), mock.Anything).
					Return([]*api.StorageObject{}, nil)
			},
		},
		{
			"null hash",
			`{"type":"core","version":"1.0.0","hash":null}`,
			nil,
			validResponse,
			func(nkMock *mockobject.NakamaModuleMock) {
				nkMock.On("StorageRead", context.Background(), mock.Anything).
					Return([]*api.StorageObject{}, nil)
				nkMock.On("StorageWrite", context.Background(), mock.Anything).
					Return([]*api.StorageObjectAck{}, nil)
			},
		},
		{
			"file does not exist",
			`{"type":"not_existed","version":"1.0.0"}`,
			errFileNotFound,
			``,
			func(nkMock *mockobject.NakamaModuleMock) {
				nkMock.On("StorageRead", context.Background(), mock.Anything).
					Return([]*api.StorageObject{}, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			loggerMock := mockobject.NewLoggerMock(t)
			nkMock := mockobject.NewNakamaModuleMock(t)

			tt.expectedCalls(nkMock)

			response, err := RpcTransferFile(ctx, loggerMock, nil, nkMock, tt.payload)

			if tt.expectedError == nil {
				require.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Equal(t, tt.expectedError, err)
			}

			assert.Equal(t, response, tt.expectedResponse)
		})
	}
}

func TestValidate(t *testing.T) {
	validRequest := RpcTransferFileRequest{
		Type:    "core",
		Version: "1.0.0",
	}
	err := validRequest.Validate()
	require.NoError(t, err)

	invalidRequest := RpcTransferFileRequest{
		Type:    "invalid/type",
		Version: "1.0.0",
	}
	err = invalidRequest.Validate()
	require.Error(t, err)

	invalidVersionRequest := RpcTransferFileRequest{
		Type:    "core",
		Version: "invalid_version/",
	}
	err = invalidVersionRequest.Validate()
	require.Error(t, err)
}
