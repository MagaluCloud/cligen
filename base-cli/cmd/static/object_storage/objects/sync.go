package objects

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type syncOptions struct {
	Local     string
	Bucket    string
	Delete    bool
	BatchSize int
}

type syncResult struct {
	Source        string `json:"src"`
	Destination   string `json:"dst"`
	FilesDeleted  int    `json:"deleted"`
	FilesUploaded int    `json:"uploaded"`
	Deleted       bool   `json:"hasDeleted"`
	DeletedFiles  string `json:"deletedFiles"`
}

func SyncCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts syncOptions

	cmd := &cobra.Command{
		Use:   "sync [local] [bucket]",
		Short: manager.T("cli.auth.object_storage.objects.sync.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSync(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Local, "local", "", manager.T("cli.auth.object_storage.objects.sync.local"))
	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", manager.T("cli.auth.object_storage.objects.sync.bucket"))
	cmd.Flags().BoolVar(&opts.Delete, "delete", false, manager.T("cli.auth.object_storage.objects.sync.delete"))
	cmd.Flags().IntVar(&opts.BatchSize, "batch-size", 1000, manager.T("cli.auth.object_storage.objects.sync.batch_size"))

	return cmd
}

func runSync(ctx context.Context, objectService objSdk.ObjectService, args []string, opts syncOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	local := opts.Local

	if len(args) > 0 {
		local = args[0]
	}

	if local == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho local ou usar a flag --local")

		return nil
	}

	bucket := opts.Bucket

	if len(args) > 1 {
		bucket = args[1]
	}

	if bucket == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do bucket ou usar a flag --bucket")

		return nil
	}

	if isRemote(local) {
		beautiful.NewOutput(rawMode).PrintError("local cannot be an bucket! To copy or move between buckets, use \"mgc object-storage objects copy/move\"")

		return nil
	}

	bucketName, _ := common.ParseBucketNameAndObjectKey(bucket)

	fileInfo, err := os.Stat(local)
	if err != nil {
		return fmt.Errorf("erro ao acessar o diretório local: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("o caminho local deve ser um diretório")
	}

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	localFiles := make(map[string]string)
	if err := collectLocalFiles(local, "", localFiles); err != nil {
		return fmt.Errorf("erro ao coletar arquivos locais: %w", err)
	}

	remoteFiles := make(map[string]bool)
	if err := collectRemoteFiles(ctx, objectService, bucketName, "", remoteFiles); err != nil {
		return fmt.Errorf("erro ao listar objetos no bucket: %w", err)
	}

	filesToUpload := []string{}
	filesToDelete := []string{}
	filesUploaded := 0
	filesDeleted := 0

	for localPath := range localFiles {
		if !remoteFiles[localPath] {
			filesToUpload = append(filesToUpload, localPath)
		}
	}

	if opts.Delete {
		for remotePath := range remoteFiles {
			if _, exists := localFiles[remotePath]; !exists {
				filesToDelete = append(filesToDelete, remotePath)
			}
		}
	}

	for i := 0; i < len(filesToUpload); i += opts.BatchSize {
		end := min(i+opts.BatchSize, len(filesToUpload))

		batch := filesToUpload[i:end]
		for _, filePath := range batch {
			if err := uploadFile(ctx, objectService, bucketName, local, filePath); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao fazer upload de %s: %v\n", filePath, err)
				continue
			}
			filesUploaded++
		}
	}

	for i := 0; i < len(filesToDelete); i += opts.BatchSize {
		end := min(i+opts.BatchSize, len(filesToDelete))

		batch := filesToDelete[i:end]
		for _, filePath := range batch {
			if err := objectService.Delete(ctx, bucketName, filePath, nil); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao deletar %s: %v\n", filePath, err)
				continue
			}
			filesDeleted++
		}
	}

	result := syncResult{
		Source:        local,
		Destination:   bucket,
		FilesDeleted:  filesDeleted,
		FilesUploaded: filesUploaded,
		Deleted:       len(filesToDelete) > 0,
		DeletedFiles:  strings.Join(filesToDelete, ", "),
	}

	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}

func collectLocalFiles(basePath, prefix string, files map[string]string) error {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(basePath, entry.Name())
		relativePath := filepath.Join(prefix, entry.Name())
		relativePath = filepath.ToSlash(relativePath)

		if entry.IsDir() {
			if err := collectLocalFiles(fullPath, relativePath, files); err != nil {
				return err
			}
		} else {
			hash, err := calculateFileMD5(fullPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao calcular MD5 de %s: %v\n", fullPath, err)
				continue
			}
			files[relativePath] = hash
		}
	}

	return nil
}

func collectRemoteFiles(ctx context.Context, objectService objSdk.ObjectService, bucketName, prefix string, files map[string]bool) error {
	objects, err := objectService.List(ctx, bucketName, objSdk.ObjectListOptions{
		Prefix: "",
	})
	if err != nil {
		return err
	}

	for _, obj := range objects {
		files[obj.Key] = true
	}

	return nil
}

func uploadFile(ctx context.Context, objectService objSdk.ObjectService, bucketName, basePath, relativePath string) error {
	fullPath := filepath.Join(basePath, filepath.FromSlash(relativePath))

	fileBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	contentType := "application/octet-stream"
	if len(fileBytes) > 0 {
		contentType = http.DetectContentType(fileBytes)
	}

	return objectService.Upload(ctx, bucketName, relativePath, fileBytes, contentType, nil)
}

func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
