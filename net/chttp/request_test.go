package chttp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChttp_GetMethodCode(t *testing.T) {
	_, err := GetUrl(context.Background(), "http://httpbin.org//redirect-to?url=www.baidu.com")
	if assert.NotNil(t, err) {
		assert.Equal(t, "http状态码404", err.Error(), "状态码404")
	}

	_, err = GetUrl(context.Background(), "http://httpbin.org/status/200")
	assert.Nil(t, err)
}
