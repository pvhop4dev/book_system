package middleware

import (
	"book_system/i18n"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

type Error struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	I18nData map[string]string `json:"i18nData"`
}

func (e Error) GetMesssageI18n(lang string) string {
	return i18n.LocalizeWithValue(e.Message, lang, e.I18nData)
}

func (e Error) Error() string {
	return e.GetMesssageI18n(e.Message)
}

var Success = Error{
	Code:     200,
	Message:  "success",
	I18nData: map[string]string{},
}

var Forbidden = Error{
	Code:     403,
	Message:  "forbidden",
	I18nData: map[string]string{},
}

var NotFound = Error{
	Code:     404,
	Message:  "not_found",
	I18nData: map[string]string{},
}

var InternalServerError = Error{
	Code:     500,
	Message:  "internal_server_error",
	I18nData: map[string]string{},
}

var FileEmpty = Error{
	Code:     400,
	Message:  "file_empty",
	I18nData: map[string]string{},
}

var Unauthorized = Error{
	Code:     401,
	Message:  "unauthorized",
	I18nData: map[string]string{},
}

var BadRequest = Error{
	Code:     400,
	Message:  "bad_request",
	I18nData: map[string]string{},
}

var UploadFileFailed = Error{
	Code:     400,
	Message:  "upload_file_failed",
	I18nData: map[string]string{},
}

var EmailInvalid = Error{
	Code:     400,
	Message:  "email_invalid",
	I18nData: map[string]string{"email": "email"},
}

var FieldRequired = Error{
	Code:     400,
	Message:  "field_required",
	I18nData: map[string]string{"field": "field"},
}

var FieldInvalid = Error{
	Code:     400,
	Message:  "field_invalid",
	I18nData: map[string]string{"field": "field"},
}
