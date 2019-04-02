package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/alex-phillips/lychee/lib/log"
)

var DryRun bool = false

func Chunk(items []string) (retval [][]string) {
	numCPUs := runtime.NumCPU()
	chunkSize := (len(items) + numCPUs - 1) / numCPUs

	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize

		if end > len(items) {
			end = len(items)
		}

		retval = append(retval, items[i:end])
	}

	return retval
}

func HttpGet(source string, cookies map[string]string) ([]byte, error) {
	log.Debug.Println("Requesting " + source)
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		log.Error.Fatal(err)
	}

	for key, value := range cookies {
		cookie := http.Cookie{Name: key, Value: value}
		req.AddCookie(&cookie)
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36")

	resp, getErr := client.Do(req)
	if getErr != nil {
		return []byte{}, getErr
	}

	log.Debug.Println(fmt.Sprintf("Response code: %d", resp.StatusCode))
	if resp.StatusCode == 404 {
		return []byte{}, errors.New(source + " returned 404")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return body, nil
}

func DownloadFile(path string, source string) error {
	if DryRun == true {
		log.Debug.Println("Dry run - not downloading file")
		return nil
	}

	log.Debug.Println("Downloading from URL " + source)

	res, err := http.Head(source)
	if err != nil {
		panic(err)
	}
	contentlength := res.ContentLength

	// Try to get extension from response filetype
	ext := ""
	mimeExt, err := mime.ExtensionsByType(res.Header.Get("Content-Type"))
	if err != nil || len(mimeExt) == 0 {
		log.Debug.Println("Unable to determine extension from content-type " + res.Header.Get("Content-Type"))
		u, _ := url.Parse(source)
		tmpSrc := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
		log.Debug.Println("Query parameters trimmed: " + tmpSrc)
		ext = filepath.Ext(tmpSrc)
	} else {
		ext = mimeExt[0]
	}

	path = path + ext

	if f, err := os.Stat(path); err == nil {
		if f.Size() == contentlength {
			log.Debug.Println("File " + path + " exists and is the same size. Skipping.")

			return nil
		}
	}

	log.Debug.Println(fmt.Sprintf("Downloading file %s to %s", source, path))

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetJSONValue(obj interface{}, path string) interface{} {
	pathParts := strings.Split(path, ".")
	var retval interface{}
	for _, key := range pathParts {
		retval = obj.(map[string]interface{})[key]
	}

	return retval
}

func Slugify(s string) string {
	var re = regexp.MustCompile("[^a-z0-9]+")

	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
