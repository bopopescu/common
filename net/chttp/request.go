package chttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const DialTimeout = time.Second * 1
const KeepAliveTime = time.Second * 15
const HttpRequestTimeout = time.Second * 3

func _httpCodeError(code int) error {
	return errors.Errorf("http状态码%d", code)
}

// NOTE: tls == nil mean http request
func _getHttpClient(ctx context.Context, tls *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   DialTimeout,
				KeepAlive: KeepAliveTime,
			}).DialContext,
			TLSHandshakeTimeout: DialTimeout,
			TLSClientConfig:     tls,
		},
		Timeout: HttpRequestTimeout,
	}
}

func GenQueryStr(mp map[string]string) string {
	params := url.Values{}
	for k, v := range mp {
		params.Add(k, v)
	}
	return params.Encode()
}

func GetUrl(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := _getHttpClient(ctx, nil)
	rs, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		return nil, _httpCodeError(rs.StatusCode)
	}

	data, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func JsonGetUrl(ctx context.Context, url string, reply interface{}) error {
	rdata, err := GetUrl(ctx, url)
	if err != nil {
		return errors.Wrap(err, "GET请求失败")
	}

	if reply != nil {
		jerr := json.Unmarshal(rdata, reply)
		if jerr != nil {
			return errors.Wrap(jerr, "解码json数据失败")
		}
	}

	return nil
}

func PostUrl(ctx context.Context, url string, data string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte(data))
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	client := _getHttpClient(ctx, nil)
	rs, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		return nil, _httpCodeError(rs.StatusCode)
	}

	rtData, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	return rtData, nil
}

func JsonPostUrl(ctx context.Context, url string, data string, reply interface{}) error {
	rdata, err := PostUrl(ctx, url, data)
	if err != nil {
		return errors.Wrap(err, "POST请求失败")
	}

	if reply != nil {
		jerr := json.Unmarshal(rdata, reply)
		if jerr != nil {
			return errors.Wrap(jerr, "解码json数据失败")
		}
	}

	return nil
}

func TLSHttpClient(ctx context.Context, certFile, keyFile string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "加载证书失败")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return _getHttpClient(ctx, tlsConfig), nil
}

func TlsGetUrl(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	if client == nil {
		return nil, errors.New("空的client")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "新建request失败")
	}

	rs, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "GET请求失败")
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		return nil, _httpCodeError(rs.StatusCode)
	}

	rtData, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}
	return rtData, nil
}

func TlsPostUrl(ctx context.Context, client *http.Client, url string,
	data string) ([]byte, error) {
	if client == nil {
		return nil, errors.New("空的client")
	}

	buf := bytes.NewBuffer([]byte(data))
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, errors.Wrap(err, "新建request失败")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "POST请求失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, _httpCodeError(resp.StatusCode)
	}

	rtData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	return rtData, nil
}

func JsonTlsGetUrl(ctx context.Context, client *http.Client, url string, reply interface{}) error {
	rdata, err := TlsGetUrl(ctx, client, url)
	if err != nil {
		return err
	}

	if reply != nil {
		jerr := json.Unmarshal(rdata, reply)
		if jerr != nil {
			return errors.Wrap(jerr, "解码json数据失败")
		}
	}

	return nil
}

func JsonTlsPostUrl(ctx context.Context, client *http.Client, url string, data string, reply interface{}) error {
	rdata, err := TlsPostUrl(ctx, client, url, data)
	if err != nil {
		return err
	}

	if reply != nil {
		jerr := json.Unmarshal(rdata, reply)
		if jerr != nil {
			return errors.Wrap(jerr, "解码json数据失败")
		}
	}

	return nil
}
