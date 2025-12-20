package utils

import (
	"fmt"
	"regexp"
	"time"
)

// Constants moved from repository/global.go
const (
	EMAIL_MIN_LEN = 6
	EMAIL_MAX_LEN = 254

	USERNAME_MIN_LEN = 3
	USERNAME_MAX_LEN = 32

	PASSWORD_MIN_LEN = 8
	PASSWORD_MAX_LEN = 72

	TITLE_MIN_LEN = 10
	TITLE_MAX_LEN = 70

	POST_MIN_LEN = 1
	POST_MAX_LEN = 10_000

	COMMENT_MIN_LEN = 1
	COMMENT_MAX_LEN = 1_000
)

var (
	EmailExp    *regexp.Regexp
	UsernameExp *regexp.Regexp
	NameExp     *regexp.Regexp
)

func InitRegex() {
	var err error
	EmailExp, err = regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]+`)
	if err != nil {
		fmt.Println(err)
		return
	}
	UsernameExp, err = regexp.Compile(`^[a-zA-Z]+[a-zA-Z0-9._]+`)
	if err != nil {
		fmt.Println(err)
		return
	}
	NameExp, err = regexp.Compile(`^[a-zA-Z][a-zA-Z\s\-\']*$`)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ValidName(s string) bool {
	if len(s) < 1 || len(s) > 50 {
		return false
	}
	if NameExp == nil {
		InitRegex()
	}
	return NameExp.MatchString(s)
}

func ValidUsername(s string) bool {
	if len(s) < USERNAME_MIN_LEN || len(s) > USERNAME_MAX_LEN {
		return false
	}
	// Ensure regex is initialized
	if UsernameExp == nil {
		InitRegex()
	}
	return UsernameExp.MatchString(s)
}

func ValidEmail(s string) bool {
	if len(s) < EMAIL_MIN_LEN || len(s) > EMAIL_MAX_LEN {
		return false
	}
	if EmailExp == nil {
		InitRegex()
	}
	return EmailExp.MatchString(s)
}

func ValidPassword(s string) bool {
	return len(s) >= PASSWORD_MIN_LEN && len(s) <= PASSWORD_MAX_LEN
}

func ValidComment(comment string) bool {
	return len(comment) >= COMMENT_MIN_LEN && len(comment) <= COMMENT_MAX_LEN
}

func ValidPost(content string) bool {
	return len(content) >= POST_MIN_LEN && len(content) <= POST_MAX_LEN
}

func ValidPostTitle(title string) bool {
	return len(title) >= TITLE_MIN_LEN && len(title) <= TITLE_MAX_LEN
}

func SqlDateFormater(date string) string {
	today := time.Now().Format("2006-01-02")
	if len(date) != 20 {
		return "date error"
	}
	tmp := date[:10]
	if tmp == today {
		tmp = "today"
		tmp = date[11:16]
	}
	return tmp
}
