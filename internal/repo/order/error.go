package order

import (
	"net/http"

	"github.com/kargotech/gokargo/serror"
)

var ErrOrderNotFound = &serror.StdErr{
	HttpStatus: http.StatusNotFound,
	Tags:       []string{serror.TagDB, "EXPECTED"},
	ErrCode:    "REPO-ORDER-NOT-FOUND",
	ErrMsg:     "Order not found based on identifier",
}

var ErrOrderDataLayerInternal = &serror.StdErr{
	HttpStatus: http.StatusInternalServerError,
	Tags:       []string{serror.TagDB},
	ErrCode:    "REPO-ORDER-INTERNAL",
	ErrMsg:     "Internal data layer issue",
}

func init() {
	serror.RegisterErrors(
		ErrOrderNotFound,
		ErrOrderDataLayerInternal,
	)
}
