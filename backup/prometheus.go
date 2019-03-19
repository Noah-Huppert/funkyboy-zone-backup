package backup

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Noah-Huppert/golog"
	"github.com/Noah-Huppert/mountain-backup/config"
	"github.com/deckarep/golang-set"
)

// PrometheusBackuper backs up a Prometheus database. By taking a snapshot with the Prometheus admin API and backing its
// contents up under the data directory path.
type PrometheusBackuper struct {
	// Cfg for Prometheus backup job.
	Cfg config.PrometheusConfig
}

// prometheusSnapshotAPIResp is the response of the Prometheus snapshot API after a snapshot was successfully created.
type prometheusSnapshotAPIResp struct {
	// Data holds information about the snapshot.
	Data struct {
		// Name of the snapshot
		Name string `json:"name"`
	} `json:"data"`
}

// Backup Prometheus database.
func (b PrometheusBackuper) Backup(logger golog.Logger, w *tar.Writer) (int, error) {
	// {{{1 Trigger snapshot via Prometheus admin API
	// {{{2 Construct URL
	reqUrl, err := url.Parse(b.Cfg.AdminAPIHost)

	if err != nil {
		return 0, fmt.Errorf("error parsing Prometheus admin API host URL \"%s\": %s", b.Cfg.AdminAPIHost, err.Error())
	}

	reqUrl.Path = "/api/v1/admin/tsdb/snapshot"

	// {{{2 Make API request
	resp, err := http.Post(reqUrl.String(), "text/plain", bytes.NewReader([]byte{}))

	defer func() {
		if err = resp.Body.Close(); err != nil {
			logger.Fatalf("error closing Prometheus snapshot API response body: %s", err.Error())
		}
	}()

	if err != nil {
		return 0, fmt.Errorf("error making Prometheus snapshot API request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		respBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return 0, fmt.Errorf("non-OK response from Prometheus snapshot API, additionally failed to "+
				"read response body, status: %s, error: %s", resp.Status, err.Error())
		}

		return 0, fmt.Errorf("non-OK response from Prometheus snapshot API: %s", respBytes)
	}

	// {{{2 Parse API response
	decoder := json.NewDecoder(resp.Body)

	apiResp := prometheusSnapshotAPIResp{}
	if err = decoder.Decode(&apiResp); err != nil {
		return 0, fmt.Errorf("error decoding Prometheus snapshot API response into JSON: %s", err.Error())
	}

	if len(apiResp.Data.Name) == 0 {
		return 0, fmt.Errorf("name of snapshot returned by Prometheus snapshot API was empty, cannot be")
	}

	// {{{1 Find snapshot files in Prometheus data directory
	// {{{2 Ensure exists
	snapshotDirPart := fmt.Sprintf("snapshots/%s", apiResp.Data.Name)
	snapshotDir := filepath.Join(b.Cfg.DataDirectory, snapshotDirPart)

	dirInfo, err := os.Stat(snapshotDir)
	if err != nil {
		return 0, fmt.Errorf("error stat-ing Prometheus snapshot directory: %s", err.Error())
	}

	if !dirInfo.IsDir() {
		return 0, fmt.Errorf("Prometheus snapshot directory \"%s\" is not a directory", snapshotDir)
	}

	// {{{2 Get names of files
	snapshotDirSet := mapset.NewSet()
	snapshotDirSet.Add(snapshotDir)

	snapshotFiles, err := allFiles(snapshotDirSet)
	if err != nil {
		return 0, fmt.Errorf("error collecting names of snapshot files: %s", err.Error())
	}

	// {{{2 Resolve absolute paths
	absSnapshotFiles, err := absSet(snapshotFiles)
	if err != nil {
		return 0, fmt.Errorf("error resolving snapshot absolute file paths: %s", err.Error())
	}

	// {{{1 Rewrite Prometheus data paths
	// rewrittenFiles holds snapshot file paths as keys, and the path they would exist if they were in the main
	// data directory as values.
	rewrittenFiles := map[string]string{}
	absSnapshotFilesIt := absSnapshotFiles.Iterator()

	for fileUntyped := range absSnapshotFilesIt.C {
		file := fileUntyped.(string)

		rewritten := strings.ReplaceAll(file, fmt.Sprintf("/%s", snapshotDirPart), "")

		rewrittenFiles[file] = rewritten
	}

	// {{{1 Save files
	filesCount := 0

	for snapshotFile, rewrittenFile := range rewrittenFiles {
		// {{{2 Write tar header
		// {{{3 Get file info
		fileInfo, err := os.Stat(snapshotFile)
		if err != nil {
			return 0, fmt.Errorf("error stat-ing \"%s\": %s", snapshotFile, err.Error())
		}

		// {{{3 Write header
		err = w.WriteHeader(&tar.Header{
			Name: rewrittenFile,
			Mode: int64(fileInfo.Mode().Perm()),
			Size: fileInfo.Size(),
		})

		// {{{2 Write file body
		// {{{3 Open file
		fileReader, err := os.Open(snapshotFile)
		if err != nil {
			return 0, fmt.Errorf("error opening \"%s\" for reading: %s", snapshotFile, err.Error())
		}

		// {{{3 Read file body
		body, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return 0, fmt.Errorf("error reading \"%s\" file contents: %s", snapshotFile, err.Error())
		}

		// {{{3 Write to tar
		if _, err = w.Write(body); err != nil {
			return 0, fmt.Errorf("error writing \"%s\" to tar file: %s", snapshotFile, err.Error())
		}

		logger.Info(snapshotFile)
		filesCount++
	}

	return filesCount, nil
}
