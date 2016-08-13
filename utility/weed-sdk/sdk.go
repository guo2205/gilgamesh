// sdk
package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gilgamesh/utility/config"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type WeedSdk struct {
	option *config.ResourceOption
}

func NewWeedSdk(option *config.ResourceOption) *WeedSdk {
	return &WeedSdk{
		option: option,
	}
}

func (c *WeedSdk) UpdateFile(fid string, fileData []byte) (int, error) {
	url, err := c.lookup(fid)
	if err != nil {
		return -1, err
	}

	return c.upload(fid, url, fileData)
}

func (c *WeedSdk) RemoveFile(fid string) error {
	return c.remove(fid)
}

func (c *WeedSdk) SaveFile(fileData []byte) (string, int, error) {
	fid, url, err := c.assign()
	if err != nil {
		return "", -1, err
	}
	size, err := c.upload(fid, url, fileData)
	if err != nil {
		return "", -1, err
	}
	return fid, size, nil
}

func (c *WeedSdk) GetFile(fid string) ([]byte, error) {
	return c.getData(fid)
}

var (
	ErrOperatorFailed           error = errors.New("optrator failed")
	ErrHttpReqFailed            error = errors.New("http request failed")
	ErrBuildHttpReqFailed       error = errors.New("build http request failed")
	ErrParsingResponseFailed    error = errors.New("parsing response data failed")
	ErrFidFormatIncorrect       error = errors.New("fid format incorrect")
	ErrLocationsFormatIncorrect error = errors.New("locations format incorrect")
)

func (c *WeedSdk) assign() (fid string, url string, e error) {
	resp, err := http.Get(fmt.Sprintf(`http://%s/dir/assign`, c.option.StoreServer))
	if err != nil {
		return "", "", ErrHttpReqFailed
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", ErrHttpReqFailed
	}

	assignResult := struct {
		Fid string
		Url string
	}{}
	err = json.Unmarshal(d, &assignResult)
	if err != nil {
		return "", "", ErrParsingResponseFailed
	}

	if assignResult.Fid == "" || assignResult.Url == "" {
		return "", "", ErrOperatorFailed
	}

	return assignResult.Fid, assignResult.Url, nil
}

func (c *WeedSdk) upload(fid string, url string, fileData []byte) (size int, e error) {
	body := bytes.NewBuffer(make([]byte, 0, len(fileData)+256))
	multiWriter := multipart.NewWriter(body)
	fileWriter, _ := multiWriter.CreateFormFile("file", fid)
	fileWriter.Write(fileData)
	multiWriter.Close()

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/%s", url, fid), body)
	if err != nil {
		return -1, ErrBuildHttpReqFailed
	}
	req.Header.Add("Content-Type", multiWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, ErrHttpReqFailed
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, ErrHttpReqFailed
	}

	uploadResult := struct {
		Size int
	}{}
	err = json.Unmarshal(d, &uploadResult)
	if err != nil {
		return -1, ErrParsingResponseFailed
	}

	if uploadResult.Size == 0 {
		return -1, ErrOperatorFailed
	}

	return uploadResult.Size, nil
}

func (c *WeedSdk) lookup(fid string) (url string, e error) {
	if len(strings.Split(fid, ",")) != 2 {
		return "", ErrFidFormatIncorrect
	}

	volumeId := strings.Split(fid, ",")[0]

	resp, err := http.Get(fmt.Sprintf(`http://%s/dir/lookup?volumeId=%s`, c.option.StoreServer, volumeId))
	if err != nil {
		return "", ErrHttpReqFailed
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrHttpReqFailed
	}

	lookupResult := struct {
		Locations []struct {
			PublicUrl string
			Url       string
		}
	}{}

	err = json.Unmarshal(d, &lookupResult)
	if err != nil {
		return "", ErrParsingResponseFailed
	}

	if len(lookupResult.Locations) == 0 {
		return "", ErrLocationsFormatIncorrect
	}

	if lookupResult.Locations[0].Url == "" {
		return "", ErrOperatorFailed
	}

	return lookupResult.Locations[0].Url, nil
}

func (c *WeedSdk) getData(fid string) (data []byte, e error) {
	url, err := c.lookup(fid)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf(`http://%s/%s`, url, fid))
	if err != nil {
		return nil, ErrHttpReqFailed
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrHttpReqFailed
	}

	if len(d) == 0 {
		return nil, ErrOperatorFailed
	}

	return d, nil
}

func (c *WeedSdk) remove(fid string) (e error) {
	url, err := c.lookup(fid)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/%s", url, fid), nil)
	if err != nil {
		return ErrBuildHttpReqFailed
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ErrHttpReqFailed
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrHttpReqFailed
	}

	deleteResult := struct {
		Size int
	}{}
	err = json.Unmarshal(d, &deleteResult)
	if err != nil {
		return ErrParsingResponseFailed
	}

	if deleteResult.Size == 0 {
		return ErrOperatorFailed
	}

	return nil
}
