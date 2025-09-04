package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sopial42/bifrost/pkg/logger"
)

type ErrResponse struct {
	Error ErrDetails `json:"error"`
}

type ErrDetails struct {
	AppCode AppErrorCode `json:"app_code"`
	Message string       `json:"message"`
	TraceID string       `json:"trace_id,omitempty"`
	Origin  error        `json:"-"`
}

// type ErrDetails struct {
// 	AppError
// 	TraceID string `json:"trace_id,omitempty"`
// }

var unexpectedErrMessage = ErrResponse{
	Error: ErrDetails{
		AppCode: CodeErrUnexpected,
		Message: "internal server error",
	},
}

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
		errResponse := ErrResponse{
			Error: ErrDetails{
				AppCode: appErr.Code,
				Message: appErr.Message,
				Origin:  appErr.Origin,
			},
		}

		// if traceID := tracing.GetTracingIDFromContext(c); traceID != "" {
		// 	errMessage.Error.TraceID = traceID
		// }

		ctx.JSON(errResponse.GetHTTPCode(), errResponse)
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
			errResponse := ErrResponse{
				Error: ErrDetails{
					Message: msgStr,
				},
			}

			ctx.JSON(httpErr.Code, errResponse)
			return
		}
	}

	// Handle if its an unexpected error
	// if traceID := tracing.GetTracingIDFromContext(c); traceID != "" {
	// 	unexpectedErrMessage.TraceID = traceID
	// }

	unexpectedErrMessage.Error.Origin = err
	log.Err(err).Errorf("HTTP unexpected error")
	ctx.JSON(http.StatusInternalServerError, unexpectedErrMessage)
}

func (e *ErrResponse) GetHTTPCode() int {
	switch e.Error.AppCode {
	case CodeErrInvalidInput:
		return http.StatusBadRequest
	case CodeErrAlreadyExists:
		return http.StatusConflict
	case CodeErrUnauthorized:
		return http.StatusUnauthorized
	case CodeErrForbidden:
		return http.StatusForbidden
	case CodeErrNotFound:
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func SetCustomErrorHandler(echo *echo.Echo) {
	echo.HTTPErrorHandler = ErrorsHandler
}
