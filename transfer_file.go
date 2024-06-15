package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/heroiclabs/nakama-common/runtime"
)

// RpcTransferFileRequest represents the payload for the RPC transfer file request.
type RpcTransferFileRequest struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

// RpcTransferFileResponse represents the response for the RPC transfer file request.
type RpcTransferFileResponse struct {
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Hash    string          `json:"hash"`
	Content json.RawMessage `json:"content"`
}

// FileStorageItem represents the file storage item to be saved in Nakama's storage.
type FileStorageItem struct {
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Hash    string          `json:"hash"`
	Content json.RawMessage `json:"content"`
}

const defaultType = "core"
const defaultVersion = "1.0.0"
const filestorageCollection = "filestorage"

var validTypePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
var validVersionPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// Validate checks the fields of RpcTransferFileRequest.
func (request *RpcTransferFileRequest) Validate() (err error) {
	if !validTypePattern.MatchString(request.Type) {
		err = errors.New("invalid type field in request")
	}
	if !validVersionPattern.MatchString(request.Version) {
		err = errors.Join(err, errors.New("invalid version field in request"))
	}
	return err
}

// RpcTransferFile handles the RPC call to transfer a file. It reads the file content from the disk, calculates its hash,
// and stores the data in the Nakama storage if not already present. It returns the file content and its hash.
func RpcTransferFile(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if ok && userId != "" {
		logger.Error("rpc was called by a user")
		return "", errCtxUserIdFound
	}

	var request RpcTransferFileRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", errUnmarshal
	}

	if request.Type == "" {
		request.Type = defaultType
	}
	if request.Version == "" {
		request.Version = defaultVersion
	}

	if err := request.Validate(); err != nil {
		return "", errValidationFailed
	}

	resp := RpcTransferFileResponse{
		Type:    request.Type,
		Version: request.Version,
	}

	key := fmt.Sprintf("%s-%s", request.Type, request.Version)

	storageItem, err := findFileStorageItem(ctx, logger, db, nk, key)

	if err != nil {
		return "", err
	}

	if storageItem != nil {
		resp.Hash = storageItem.Hash

		// If the hash is empty or the hash matches the content, return the content
		if request.Hash == "" || (request.Hash != "" && request.Hash == storageItem.Hash) {
			resp.Content = storageItem.Content
		}

		responseBytes, err := json.Marshal(resp)
		if err != nil {
			return "", errMarshal
		}

		return string(responseBytes), nil
	}

	// Load the file from the disk
	filePath := filepath.Join(request.Type, fmt.Sprintf("%s.json", request.Version))
	content, err := fs.ReadFile(filestorage, filePath)
	if err != nil {
		return "", errFileNotFound
	}

	// Calculate the hash of the file content
	fileHash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(fileHash[:])

	resp.Hash = hashStr

	// If the hash is empty or the hash matches the content, return the content and store it in db
	if request.Hash == "" || (request.Hash != "" && request.Hash == hashStr) {
		resp.Content = content

		item := FileStorageItem{
			Type:    request.Type,
			Version: request.Version,
			Hash:    hashStr,
			Content: resp.Content,
		}

		err := saveFileStorageItem(ctx, logger, db, nk, item, key)
		if err != nil {
			return "", err
		}
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return "", errMarshal
	}

	return string(responseBytes), nil
}

// findFileStorageItem retrieves the file storage item from Nakama's storage using the provided key.
func findFileStorageItem(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, key string) (*FileStorageItem, error) {
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: filestorageCollection,
		Key:        key,
	}})

	if err != nil {
		logger.Error("StorageRead error: %v", err)
		return nil, errInternalError
	}

	if len(objects) == 0 {
		return nil, nil
	}

	storageItem := &FileStorageItem{}
	if err := json.Unmarshal([]byte(objects[0].GetValue()), storageItem); err != nil {
		return nil, errUnmarshal
	}

	return storageItem, nil
}

// saveFileStorageItem stores the file storage item in Nakama's storage.
func saveFileStorageItem(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, storageItem FileStorageItem, key string) error {
	storageItemEncoded, err := json.Marshal(storageItem)
	if err != nil {
		return errMarshal
	}

	// Store the data using Nakama's storage engine
	writeStorageObjects := []*runtime.StorageWrite{
		{
			Collection:      filestorageCollection,
			Key:             key,
			Value:           string(storageItemEncoded),
			PermissionRead:  runtime.STORAGE_PERMISSION_OWNER_READ,
			PermissionWrite: runtime.STORAGE_PERMISSION_OWNER_WRITE,
		},
	}

	if _, err := nk.StorageWrite(ctx, writeStorageObjects); err != nil {
		logger.Error("StorageWrite error: %v", err)
		return errInternalError
	}

	return nil
}
