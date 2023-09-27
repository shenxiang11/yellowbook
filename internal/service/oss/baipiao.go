package oss

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type IService interface {
	Upload(f *multipart.FileHeader) (string, error)
}

type Service struct {
	client *http.Client
}

func NewService() IService {
	return &Service{client: http.DefaultClient}
}

type baiPiaoResp struct {
	Code    int64  `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

func (s *Service) Upload(f *multipart.FileHeader) (string, error) {
	url := "https://front-gateway.mollybox.com/service-person-center/api/pet/uploadPetAvatar"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := f.Open()
	defer file.Close()

	part1,
		err := writer.CreateFormFile("avatar", filepath.Base(f.Filename))
	_, err = io.Copy(part1, file)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}

	req.Header.Add("Host", "front-gateway.mollybox.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("xweb_xhr", "1")
	req.Header.Add("contentType", "multipart/form-data")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF XWEB/30817")
	req.Header.Add("token", viper.GetString("baipiao_oss_token"))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://servicewechat.com/wx2c9036bbb0046e67/235/page-frame.html")
	req.Header.Add("Accept-Language", "zh-CN,zh")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resp baiPiaoResp
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return "", err
	}

	return resp.Data, nil
}
