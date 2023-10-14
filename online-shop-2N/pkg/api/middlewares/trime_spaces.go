package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"online-shop-2N/pkg/api/handlers/responses"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// TrimSpaces implements Middleware.
func (*middleware) TrimSpaces() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(32<<20))

		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Abort()
			responses.ErrorResponse(ctx, http.StatusBadRequest, "failed to ready request body", err, nil)
			return
		}

		bodyBytes = trimSpaceInJSON(bodyBytes)

		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
}

func trimSpaceInJSON(data []byte) []byte {
	var mapData map[string]interface{}

	for key, value := range mapData {
		if strValue, ok := value.(string); ok {
			mapData[key] = strings.TrimSpace(strValue)
		}
	}

	trimmedData, err := json.Marshal(mapData)
	if err != nil {
		log.Error("Failed to trim data from JSON. Due to error: ", err)
		return data
	}
	return trimmedData
}
