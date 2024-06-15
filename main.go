package main

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
	"io/fs"
	"os"
)

const (
	rpcTransferFile = "transfer_file"
)

var filestorage fs.FS

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	env := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)
	filestoragePath := env["filestorage_path"]
	filestorage = os.DirFS(filestoragePath)

	err := initializer.RegisterRpc(rpcTransferFile, RpcTransferFile)
	if err != nil {
		return err
	}
	return nil
}
