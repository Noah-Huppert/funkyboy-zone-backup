package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Noah-Huppert/mountain-backup/backup"
	"github.com/Noah-Huppert/mountain-backup/config"

	"github.com/Noah-Huppert/goconf"
	"github.com/Noah-Huppert/golog"
	"github.com/jehiah/go-strftime"
	"github.com/minio/minio-go"
	"github.com/thecodeteam/goodbye"
)

func main() {
	// {{{1 Setup goodbye library
	ctx := context.Background()
	defer goodbye.Exit(ctx, -1)

	// {{{1 Setup log
	logger := golog.NewStdLogger("backup")

	// {{{1 Load configuration
	cfgLoader := goconf.NewDefaultLoader()

	cfgLoader.AddConfigPath("./*.toml")
	cfgLoader.AddConfigPath("/etc/funkyboy-zone-backup/*.toml")

	cfg := config.Config{}
	if err := cfgLoader.Load(&cfg); err != nil {
		logger.Fatalf("error loading configuration: %s", err.Error())
	}

	// {{{1 Publish metrics on exit
	// backupSuccess will be set to false if backup failed before program exits
	backupSuccess := true

	// backupNumberFiles is the number of files which were backed up
	backupNumberFiles := 0

	goodbye.Register(func(ctx context.Context, sig os.Signal) {
		// curl --fail --silent --show-error --data-binary @- "$push_srv/metrics/job/$job"
		// {{{2 Construct request URL
		reqUrl, err := url.Parse(cfg.Metrics.PushGatewayHost)
		reqUrl.Path = fmt.Sprintf("/metrics/job/backup/host/%s", cfg.Metrics.LabelHost)

		// {{{2 Construct body
		backupSuccessInt := 1
		if !backupSuccess {
			backupSuccessInt = 0
		}
		bodyStr := fmt.Sprintf("backup_success %d\nbackup_number_files %d\n", backupSuccessInt, backupNumberFiles)
		bodyBytes := bytes.NewReader([]byte(bodyStr))

		// {{{2 Make request
		resp, err := http.Post(reqUrl.String(), "text/plain", bodyBytes)
		if err != nil {
			logger.Fatalf("error pushing metrics to Prometheus Push Gateway: %s", err.Error())
		}

		if resp.StatusCode != http.StatusAccepted {
			logger.Error("error pushing metrics to Prometheus Push Gateway, received non OK "+
				"response, status: %s", resp.Status)

			errBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Fatalf("error reading Prometheus Push Gateway response body: %s", err.Error())
			}

			logger.Fatalf("response body: %s", errBytes)
		}
	})

	// {{{1 Open tar gz file
	// {{{2 Open tar file
	backupName := fmt.Sprintf("%s.tar.gz", strftime.Format(cfg.Upload.Format, time.Now()))
	fName := fmt.Sprintf("/var/tmp/%s", backupName)

	tarGzF, err := os.Create(fName)

	if err != nil {
		logger.Fatalf("error creating tar.gz file \"%s\": %s", fName, err.Error())
	}

	defer func() {
		if err = tarGzF.Close(); err != nil {
			logger.Fatalf("error closing tar.gz file \"%s\": %s", fName, err.Error())
		}
	}()
	defer func() {
		if err = os.Remove(fName); err != nil {
			logger.Fatalf("error removing tar.gz file \"%s\": %s", fName, err.Error())
		}
	}()

	// {{{2 Create gzip writer
	gzipW := gzip.NewWriter(tarGzF)

	defer func() {
		if err := gzipW.Close(); err != nil {
			logger.Fatalf("error closing gzip writer \"%s\": %s", fName, err.Error())
		}
	}()

	// {{{2 Create tar writer
	tarW := tar.NewWriter(gzipW)

	defer func() {
		if err := tarW.Close(); err != nil {
			logger.Fatalf("error closing tar writer \"%s\": %s", fName, err.Error())
		}
	}()

	// {{{1 Perform straight forward backup of files
	for key, c := range cfg.Files {
		backuperLogger := logger.GetChild(fmt.Sprintf("File.%s", key))

		backuperLogger.Infof("backing up File.%s", key)

		b := backup.FilesBackuper{
			Cfg: c,
		}

		numBackedUp, err := b.Backup(backuperLogger, tarW)
		if err != nil {
			backupSuccess = false
			logger.Fatalf("error running file backup for \"%s\": %s", key, err.Error())
		}

		backupNumberFiles += numBackedUp
	}

	// {{{1 Upload backup to s3 compatible object storage service
	logger.Infof("uploading %s", fName)

	// {{{2 Initialize minio client
	minioClient, err := minio.New(cfg.Upload.Endpoint, cfg.Upload.KeyID, cfg.Upload.SecretAccessKey, true)
	if err != nil {
		logger.Fatalf("error initializing s3 compatible object storage API client: %s", err.Error())
	}

	// {{{2 Upload
	uploadOpts := minio.PutObjectOptions{
		ContentType: "application/x-tar",
	}

	_, err = minioClient.FPutObject(cfg.Upload.Bucket, backupName, fName, uploadOpts)
	if err != nil {
		logger.Fatalf("error uploading backup: %s", err.Error())
	}
}
