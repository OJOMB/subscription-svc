package app

import (
	"fmt"
	"net/http"
	"strings"
)

type appErr struct {
	msg  string
	code int
}

const (
	svcErrInvalidInputData       = "bad input data"
	svcErrResourceNotFound       = "resource not found"
	svcErrEncounteredSystemError = "encountered system error"
	svcErrResourceConflict       = "resource state conflict"
)

var svcErrMessagesToStatusCodes = map[string]int{
	svcErrInvalidInputData:       http.StatusBadRequest,
	svcErrResourceNotFound:       http.StatusNotFound,
	svcErrEncounteredSystemError: http.StatusInternalServerError,
	svcErrResourceConflict:       http.StatusConflict,
}

func newAppErr(errMsg string, code int) *appErr {
	return &appErr{msg: errMsg, code: code}
}

func (app *App) newAppErrFromSvcErr(svcErr error) *appErr {
	svcErrSplit := strings.Split(svcErr.Error(), "-")
	svcErrPrefix := strings.Trim(svcErrSplit[0], " ")

	code, ok := svcErrMessagesToStatusCodes[svcErrPrefix]
	if !ok || len(svcErrSplit) < 2 {
		app.logger.Warnf("encountered incorrectly formatted Service Err: %v", svcErr)
		return &appErr{
			msg:  svcErr.Error(),
			code: http.StatusInternalServerError,
		}
	}

	return &appErr{
		msg:  svcErr.Error(),
		code: code,
	}
}

func (apperr *appErr) Error() string {
	return fmt.Sprintf(`{"error": "%s"}`, apperr.msg)
}

func (apperr *appErr) Code() int {
	return apperr.code
}
