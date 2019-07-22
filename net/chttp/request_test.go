package chttp

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChttp_GetMethodCode(t *testing.T) {
	//scUrl := "curl -I http://httpbin.org/status/418"
	errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(404)
	}))
	defer errServer.Close()

	_, err := GetUrl(context.Background(), errServer.URL)
	if assert.NotNil(t, err) {
		assert.Equal(t, "http状态码404", err.Error(), "状态码404")
	}

	_, err = PostUrl(context.Background(), errServer.URL, "")
	if assert.NotNil(t, err) {
		assert.Equal(t, "http状态码404", err.Error(), "状态码404")
	}

	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
	}))
	defer okServer.Close()

	_, err = GetUrl(context.Background(), okServer.URL)
	assert.Nil(t, err)
}

func TestChttpGetData(t *testing.T) {
	data := "test ok"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, data)
	}))
	defer server.Close()

	resData, err := GetUrl(context.Background(), server.URL)
	if assert.Nil(t, err) {
		assert.Equal(t, data, string(resData))
	}
}

func TestChttpPostData(t *testing.T) {
	data := "test post url"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		rdata, _ := ioutil.ReadAll(r.Body)
		log.Println(string(rdata))
		io.WriteString(w, string(rdata))
	}))
	defer server.Close()

	resData, err := PostUrl(context.Background(), server.URL, data)
	if assert.Nil(t, err) {
		assert.Equal(t, data, string(resData))
	}
}

type TestData struct {
	Val int
}

func TestChttp_JsonGetUrl(t *testing.T) {
	data := TestData{Val: 1234}
	dataStr, _ := json.Marshal(data)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, string(dataStr))
	}))
	defer server.Close()

	// 1. test normal
	output := TestData{}
	err := JsonGetUrl(context.Background(), server.URL, &output)
	if assert.Nil(t, err) {
		assert.Equal(t, data.Val, output.Val)
	}

	// 2. test nil
	err = JsonGetUrl(context.Background(), server.URL, nil)
	assert.Nil(t, err)
}

func TestChttp_JsonPostUrl(t *testing.T) {
	data := TestData{Val: 1234}
	dataStr, _ := json.Marshal(data)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		rdata, _ := ioutil.ReadAll(r.Body)
		log.Println(string(rdata))
		io.WriteString(w, string(rdata))
	}))
	defer server.Close()

	// 1. test normal
	output := TestData{}
	err := JsonPostUrl(context.Background(), server.URL, string(dataStr), &output)
	if assert.Nil(t, err) {
		assert.Equal(t, data.Val, output.Val)
	}

	// 2. test ok
	err = JsonPostUrl(context.Background(), server.URL, string(dataStr), nil)
	assert.Nil(t, err)
}

func TestChttp_TLSHttpClient(t *testing.T) {
	_, err := TlsHttpClient(context.Background(), "./test.cert", "./test.key", true)
	assert.Equal(t, nil, err, "意外收到错误: %s", err)

	_, err = TlsHttpClient(context.Background(), "not_valid", "not_valid", true)
	assert.NotNil(t, err)
}

func TestChttp_TlsGetUrl(t *testing.T) {
	data := "test ok"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, data)
	}))
	defer server.Close()

	cli, err := TlsHttpClient(context.Background(), "./test.cert", "./test.key", false)
	require.Nil(t, err)

	resData, err := TlsGetUrl(context.Background(), cli, server.URL)
	if assert.Nil(t, err) {
		assert.Equal(t, data, string(resData))
	}
}

func TestChttp_TlsPostUrl(t *testing.T) {
	data := "test post url"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		rdata, _ := ioutil.ReadAll(r.Body)
		log.Println(string(rdata))
		io.WriteString(w, string(rdata))
	}))
	defer server.Close()

	cli, err := TlsHttpClient(context.Background(), "./test.cert", "./test.key", false)
	require.Nil(t, err)

	resData, err := TlsPostUrl(context.Background(), cli, server.URL, data)
	if assert.Nil(t, err) {
		assert.Equal(t, data, string(resData))
	}
}

func TestChttp_JsonTlsGetUrl(t *testing.T) {
	data := TestData{Val: 1234}
	dataStr, _ := json.Marshal(data)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, string(dataStr))
	}))
	defer server.Close()

	cli, err := TlsHttpClient(context.Background(), "./test.cert", "./test.key", false)
	require.Nil(t, err)

	// 1. test normal
	output := TestData{}
	err = JsonTlsGetUrl(context.Background(), cli, server.URL, &output)
	if assert.Nil(t, err) {
		assert.Equal(t, data.Val, output.Val)
	}

	// 2. test nil
	err = JsonTlsGetUrl(context.Background(), cli, server.URL, nil)
	assert.Nil(t, err)
}

func TestChttp_JsonTlsPostUrl(t *testing.T) {
	data := TestData{Val: 1234}
	dataStr, _ := json.Marshal(data)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
		rdata, _ := ioutil.ReadAll(r.Body)
		log.Println(string(rdata))
		io.WriteString(w, string(rdata))
	}))
	defer server.Close()

	cli, err := TlsHttpClient(context.Background(), "./test.cert", "./test.key", false)
	require.Nil(t, err)

	// 1. test normal
	output := TestData{}
	err = JsonTlsPostUrl(context.Background(), cli, server.URL, string(dataStr), &output)
	if assert.Nil(t, err) {
		assert.Equal(t, data.Val, output.Val)
	}

	// 2. test ok
	err = JsonTlsPostUrl(context.Background(), cli, server.URL, string(dataStr), nil)
	assert.Nil(t, err)
}
