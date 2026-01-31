package load

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"client/internal/app"
)

const downloadDirName = "KeyStorageDownloads"

func DownloadFileByID(app *app.Ctx, id int64) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("invalid id: %d", id)
	}

	url := fmt.Sprintf("http://127.0.0.1:8080/file/download/%d", id)

	req := app.HTTP.R().
		SetDoNotParseResponse(true) // do not store file in RAM

	token := strings.TrimSpace(app.GetToken())
	if token != "" {
		req.SetHeader("Authorization", token)
	}

	resp, err := req.Get(url)
	if err != nil {
		return "", err
	}

	body := resp.RawBody()
	defer func() {
		if err := body.Close(); err != nil {
			log.Fatalf("Failed download file")
		}
	}()

	if resp.StatusCode() != 200 {
		b, _ := io.ReadAll(body)
		return "", fmt.Errorf("GET %s failed: status=%d body=%s", url, resp.StatusCode(), string(b))
	}

	// get file name from content disposition
	filename := filenameFromContentDisposition(resp.Header().Get("Content-Disposition"))
	if filename == "" {
		filename = "file-" + strconv.FormatInt(id, 10)
	}

	// create download directory
	baseDir, err := executableDir()
	if err != nil {
		wd, e := os.Getwd()
		if e != nil {
			return "", fmt.Errorf("cannot determine executable dir: %v (getwd failed: %v)", err, e)
		}
		baseDir = wd
	}

	outDir := filepath.Join(baseDir, downloadDirName)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", outDir, err)
	}

	filename = sanitizeFilename(filename)
	outPath := filepath.Join(outDir, filename)
	outPath = dedupePath(outPath)

	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("open file %s: %w", outPath, err)
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(f, body); err != nil {
		_ = os.Remove(outPath)
		return "", fmt.Errorf("save file: %w", err)
	}

	return outPath, nil
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, _ = filepath.EvalSymlinks(exe)
	return filepath.Dir(exe), nil
}

func filenameFromContentDisposition(v string) string {
	// attachment; filename="my.txt"
	re := regexp.MustCompile(`(?i)filename="([^"]+)"`)
	m := re.FindStringSubmatch(v)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\x00", "")
	name = filepath.Base(name)
	bad := regexp.MustCompile(`[<>:"/\\|?*\r\n\t]`)
	name = bad.ReplaceAllString(name, "_")
	if name == "" || name == "." || name == ".." {
		name = "downloaded_file"
	}
	return name
}

func dedupePath(path string) string {
	if _, err := os.Stat(path); err != nil {
		return path
	}
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for i := 1; i < 10_000; i++ {
		p := fmt.Sprintf("%s (%d)%s", base, i, ext)
		if _, err := os.Stat(p); err != nil {
			return p
		}
	}
	return fmt.Sprintf("%s (%d)%s", base, os.Getpid(), ext)
}
