package lychee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Api struct {
	Client  *http.Client
	Host    string
	version int
}

func GetApi(username string, password string, host string, version int) *Api {
	jar, _ := cookiejar.New(nil)
	api := Api{
		Client: &http.Client{Jar: jar},
		Host:   host,
	}

	api.version = version

	api.Login(username, password)

	return &api
}

func (app *Api) Login(username string, password string) bool {
	function := "Session::login"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function": {"Session::login"},
		"user":     {username},
		"password": {password},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) == "false" {
		return false
	}

	return true
}

func (app *Api) CreateAlbum(name string, parent string) string {
	function := "Album::add"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function":  {"Album::add"},
		"title":     {name},
		"parent_id": {parent},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	return string(body)
}

func (app *Api) GetAlbums() []Album {
	function := "Albums::get"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function": {function},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var results AlbumsGetResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatalln(err)
	}

	return results.Albums
}

func (app *Api) GetAlbumInfo(id string) Album {
	function := "Album::get"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function": {"Album::get"},
		"albumID":  {id},
		"password": {""},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var album Album
	err = json.Unmarshal(body, &album)
	if err != nil {
		log.Fatalln(err)
	}

	return album
}

func (app *Api) Upload(path string, albumID string) (string, error) {
	file, _ := os.Open(path)
	defer file.Close()

	formBody := &bytes.Buffer{}
	writer := multipart.NewWriter(formBody)
	_ = writer.WriteField("function", "Photo::add")
	_ = writer.WriteField("albumID", albumID)
	part, _ := writer.CreateFormFile("0", filepath.Base(file.Name()))
	io.Copy(part, file)

	writer.Close()

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/Photo::add"
	}

	req, err := http.NewRequest("POST", app.Host+requestURI, formBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := app.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	retval := string(body)
	_, err = strconv.Atoi(retval)
	if err != nil {
		return "", errors.New(retval)
	}

	return retval, nil
}

func (app *Api) Search(query string) Search {
	function := "search"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function": {function},
		"term":     {query},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var results Search
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatalln(err)
	}

	return results
}

func (app *Api) GetAlbumByName(name string) *Album {
	albums := app.GetAlbums()

	for _, album := range albums {
		if album.Title == name {
			return &album
		}
	}

	return nil
}

func (app *Api) MoveImages(ids []string, albumID string) {
	function := "Photo::setAlbum"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	var idStrs []string
	for _, id := range ids {
		idStrs = append(idStrs, id)
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function": {function},
		"photoIDs": {strings.Join(idStrs, ",")},
		"albumID":  {albumID},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
}

func (app *Api) SetTags(photoIDs []string, tags []string) {
	function := "Photo::setTags"

	requestURI := "/php/index.php"
	if app.version != 1 {
		requestURI = "/api/" + function
	}

	resp, err := app.Client.PostForm(app.Host+requestURI, url.Values{
		"function": {function},
		"photoIDs": {strings.Join(photoIDs, ",")},
		"tags":     {strings.Join(tags, ",")},
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
}
