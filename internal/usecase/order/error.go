package order

import (
	"net/http"

	"github.com/kargotech/gokargo/serror"
)

var ErrOrderValidationNotPassed = &serror.StdErr{
	HttpStatus: http.StatusBadRequest,
	ErrCode:    "UC-ORDER-VALIDATION-NOT-PASSED",
	ErrMsg:     "Validation do not pass",
	Tags:       []string{serror.TagValidation},
}

var ErrOrderInternal = &serror.StdErr{
	HttpStatus: http.StatusInternalServerError,
	ErrCode:    "UC-ORDER-INTERNAL",
	ErrMsg:     "Internal order issue",
	Tags:       []string{serror.TagDB},
}

var ErrOrderNotFound = &serror.StdErr{
	HttpStatus: http.StatusNotFound,
	ErrCode:    "UC-ORDER-NOT-FOUND",
	ErrMsg:     "Order with given argument is not found",
	Tags:       []string{serror.TagDB},
}

func init() {
	serror.RegisterErrors(
		ErrOrderValidationNotPassed,
		ErrOrderInternal,
	)
}
