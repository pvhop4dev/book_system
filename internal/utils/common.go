package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Map[T, U any](data []T, f func(T) U) []U {
	res := make([]U, 0, len(data))
	for _, e := range data {
		res = append(res, f(e))
	}
	return res
}

func Pointer[T any](val T) *T {
	return &val
}

func StrToInt(val string, defaultVal int) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func SliceStrToInt(val []string) []int {
	return Map(val, func(s string) int {
		return StrToInt(s, 0)
	})
}

func Contains(array []string, value string) bool {
	set := Slice2Set(array)
	return set[value]
}

func Slice2Set(array []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range array {
		set[v] = true
	}
	return set
}

func AnyContains(array []string, value []string) bool {
	set := Slice2Set(array)
	for _, v := range value {
		if set[v] {
			return true
		}
	}
	return false
}

func Filter[T any](data []T, f func(T) bool) []T {
	arr := make([]T, 0, len(data))
	for _, e := range data {
		if f(e) {
			arr = append(arr, e)
		}
	}
	return arr
}

func BuildTree[T any](slice []T, isParent func(m T) bool, compare func(p T, c T) bool, setTree func(p *T, c []T)) []T {
	tree := Filter(slice, isParent)
	for i := range tree {
		buildTree(&tree[i], slice, compare, setTree)
	}
	return tree
}

func buildTree[T any](parent *T, slice []T, compare func(p T, c T) bool, setTree func(p *T, c []T)) {
	var children = Filter(slice, func(c T) bool {
		return compare(*parent, c)
	})
	if len(children) == 0 {
		return
	}
	for i := range children {
		buildTree(&children[i], slice, compare, setTree)
	}
	setTree(parent, children)
}

func PageSize(c *gin.Context) (int, int) {
	page := c.Query("page")
	size := c.Query("size")
	return StrToInt(page, 0), StrToInt(size, 10)
}

func DateTimeFormat(t time.Time) string {
	return t.Format("02/01/2006 15:04:05")
}

func DateFormat(t time.Time) string {
	return t.Format("02/01/2006")
}

func TrunDate(t *time.Time, toStart bool) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}
	if toStart {
		return Pointer(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
	}
	return Pointer(time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()))
}

func GetCurrentUsername(c *gin.Context) string {
	username := c.Request.Header["Username"]
	if len(username) > 0 {
		return username[0]
	}
	return ""
}

func GetCurrentUserId(c *gin.Context) string {
	id := c.Request.Header["Id"]
	if len(id) > 0 {
		return id[0]
	}
	return ""
}

func GetCurrentLang(c *gin.Context) string {
	id := c.Request.Header["Accept-Language"]
	if len(id) > 0 {
		return id[0]
	}
	return "en"
}

func GetCurrentPosition(c *gin.Context) int {
	id := c.Request.Header["Position"]
	if len(id) > 0 {
		return StrToInt(id[0], 0)
	}
	return 0
}

func GetCurrentPartnerId(c *gin.Context) string {
	id := c.Request.Header["Partner-Id"]
	if len(id) > 0 {
		return id[0]
	}
	return ""
}

func GetCurrentProfileId(c *gin.Context) string {
	id := c.Request.Header["Profile-Id"]
	if len(id) > 0 {
		return id[0]
	}
	return ""
}

/*
*
Dùng cho lưu audit log
Các trường cần lưu thì thêm tag audit
Các value cần transform thì thêm tag audit_transform
*/

func GeneratePassword(length int) string {
	var (
		lowerCharSet    = "abcdedfghijklmnopqrstuvwxyz"
		upperCharSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numberSet       = "0123456789"
		specicalCharSet = "!@#$%&*"
		allCharSet      = lowerCharSet + upperCharSet + numberSet
	)

	var password strings.Builder
	minNum := 1
	minUpperCase := 1
	minLowerCase := 1
	minSpecicalCase := 1

	for i := 0; i < minLowerCase; i++ {
		random := rand.Intn(len(lowerCharSet))
		password.WriteString(string(lowerCharSet[random]))
	}

	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	for i := 0; i < minSpecicalCase; i++ {
		random := rand.Intn(len(specicalCharSet))
		password.WriteString(string(specicalCharSet[random]))
	}

	remainingLength := length - minLowerCase - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
