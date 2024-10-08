package xhttp

import "net/http"

const (
	STATUS_OK                    = http.StatusOK
	STATUS_CREATED               = http.StatusCreated
	STATUS_BAD_REQUEST           = http.StatusBadRequest
	STATUS_INTERNAL_SERVER_ERROR = http.StatusInternalServerError
	STATUS_SEE_OTHER             = http.StatusSeeOther
)

func Text(status int) string {
	return http.StatusText(status)
}
