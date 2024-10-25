package middleware

import (
	"net/http"

	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models/external"
	"github.com/gin-gonic/gin"
)

var GetOrderListReqParams = map[string]bool{
	"limit":  true,
	"offset": true,
}

var AllowedQueryParams = map[string]map[string]bool{
	http.MethodGet + "/ecommerce/v1/orders":        GetOrderListReqParams,
	http.MethodPost + "/ecommerce/v1/orders":       nil,
	http.MethodGet + "/ecommerce/v1/orders/:id":    nil,
	http.MethodDelete + "/ecommerce/v1/orders/:id": nil,
}

func QueryParamsCheckMiddleware(lgr *logger.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, requestID := lgr.WithReqID(c)

		allowedQueryParams, ok := AllowedQueryParams[c.Request.Method+c.FullPath()]
		if !ok {
			l.Error().
				Str("method", c.Request.Method).
				Str("path", c.FullPath()).
				Msg("unspuuorted method or path")
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusNotFound,
				ErrorCode:      "",
				Message:        "Unsupported method or path",
				DebugID:        requestID,
			}
			c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
			return
		}

		hasBadReqParams := HasUnSupportedQueryParams(c.Request, allowedQueryParams)
		if hasBadReqParams {
			l.Error().Str("given query params", c.Request.URL.RawQuery).
				Interface("allowed query params", allowedQueryParams).
				Str("requestPath", c.FullPath()).
				Str("requestMethod", c.Request.Method).
				Msg("request hash unsupported query params")

			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        "invalid query params",
				DebugID:        requestID,
			}
			c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		c.Next()

	}

}

func HasUnSupportedQueryParams(req *http.Request, supportedParams map[string]bool) bool {
	queryParams := req.URL.Query()
	for param := range queryParams {
		if _, ok := supportedParams[param]; !ok {
			return true

		}
	}
	return false
}
