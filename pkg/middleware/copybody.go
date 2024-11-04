package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/supermicah/go-framework-admin/pkg/errors"
	"github.com/supermicah/go-framework-admin/pkg/util"
)

type CopyBodyConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	MaxContentLen       int64
}

var DefaultCopyBodyConfig = CopyBodyConfig{
	MaxContentLen: 32 << 20, // 32MB
}

func CopyBody() gin.HandlerFunc {
	return CopyBodyWithConfig(DefaultCopyBodyConfig)
}

func CopyBodyWithConfig(config CopyBodyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			c.Request.Body == nil {
			c.Next()
			return
		}

		var (
			requestBody []byte
			err         error
		)

		isGzip := false
		safe := http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxContentLen)
		if c.GetHeader("Content-Encoding") == "gzip" {
			if reader, errr := gzip.NewReader(safe); errr == nil {
				isGzip = true
				requestBody, err = io.ReadAll(reader)
			}
		}

		if !isGzip {
			requestBody, err = io.ReadAll(safe)
		}

		if err != nil {
			util.ResError(c, errors.RequestEntityTooLarge("", "Request body too large, limit %d byte", config.MaxContentLen))
			return
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = io.NopCloser(bf)
		c.Set(util.ReqBodyKey, requestBody)
		c.Next()
	}
}
