package errors

import (
	"errors"
	"net/http"

	"github.com/bifrost/internal/common/logger"
	"github.com/labstack/echo/v4"
)

var unexpectedErrMessage = map[string]string{"error": "internal server error"}

//nolint:errcheck
var ErrorsHandler = func(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	// Handle appErrors
	var appErr *AppError
	c := ctx.Request().Context()
	log := logger.GetLogger(c)

	if errors.As(err, &appErr) {
		log.Debugf("AppError: %v", appErr)
		errMessage := map[string]string{
			"error": err.Error(),
		}

		switch appErr.Code {
		case ErrInvalidInput:
			ctx.JSON(http.StatusBadRequest, errMessage)
		case ErrAlreadyExists:
			ctx.JSON(http.StatusConflict, errMessage)
		case ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, errMessage)
		case ErrForbidden:
			ctx.JSON(http.StatusForbidden, errMessage)
		case ErrNotFound:
			ctx.JSON(http.StatusNotFound, errMessage)
		default:
			log.Err(appErr).Errorf("AppError unhandled")
			// if traceID := tracing.GetTracingIDFromContext(c); traceID != "" {
			// 	unexpectedErrMessage["trace_id"] = traceID
			// }
			ctx.JSON(http.StatusInternalServerError, unexpectedErrMessage)
		}

		return
	}

	// Handle if its a client error
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.Code < http.StatusInternalServerError {
			message := httpErr.Message
			msgStr, ok := message.(string)
			if !ok {
				msgStr = httpErr.Error()
			}

			log.Err(err).Debugf("HTTP error")
			errMessage := map[string]string{
				"error": msgStr,
			}

			ctx.JSON(httpErr.Code, errMessage)
			return
		}
	}

	// Handle if its an unexpected error
	// if traceID := tracing.GetTracingIDFromContext(c); traceID != "" {
	// 	unexpectedErrMessage["trace_id"] = traceID
	// }

	log.Err(err).Errorf("HTTP unexpected error")
	ctx.JSON(http.StatusInternalServerError, unexpectedErrMessage)
}

func SetCustomErrorHandler(echo *echo.Echo) {
	echo.HTTPErrorHandler = ErrorsHandler
}
