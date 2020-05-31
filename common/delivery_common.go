package common

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const errorFormParse = "cannot parse %s from request form"

// ParsePaginationForm is used to Parse Pagination Form from request
func ParsePaginationForm(r *http.Request) (Pagination, error) {
	errorMsgs := make([]string, 0)
	page, err := strconv.ParseInt(r.FormValue("p"), 10, 64)
	if err != nil {
		errorMsgs = append(errorMsgs, fmt.Sprintf(errorFormParse, "page"))
		err = nil
	}

	limit, err := strconv.ParseInt(r.FormValue("l"), 10, 64)
	if err != nil {
		errorMsgs = append(errorMsgs, fmt.Sprintf(errorFormParse, "limit"))
		err = nil
	}

	if len(errorMsgs) > 0 {
		return Pagination{}, fmt.Errorf(strings.Join(errorMsgs, `\n`))
	}

	return Pagination{
		Page:  int(page),
		Limit: int(limit),
	}, nil
}

// Pagination is the structure data for pagination
type Pagination struct {
	Page  int
	Limit int
}
