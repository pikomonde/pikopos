package common

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const salt = "ymOcixH29JWua0684iEd"

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s + salt))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func OTP(digit int) (string, error) {
	b := make([]byte, digit)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = '0' + byte(math.Floor(float64(b[i])/25.6))
	}
	return string(b), nil
}

// DBTx is used as an interface for sqlx.Tx and sqlx.DB
type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

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
