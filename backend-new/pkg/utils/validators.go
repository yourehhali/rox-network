package utils

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

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
	POST_MAX_LEN = 5000

	COMMENT_MIN_LEN = 1
	COMMENT_MAX_LEN = 500

	NAME_MIN_LEN = 1
	NAME_MAX_LEN = 50

	CHAT_MSG_MIN_LEN = 1
	CHAT_MSG_MAX_LEN = 250

	GROUP_NAME_MIN_LEN = 3
	GROUP_NAME_MAX_LEN = 50

	GROUP_DESC_MAX_LEN = 500

	ABOUT_MAX_LEN = 500
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
		log.Println(err)
		return
	}
	UsernameExp, err = regexp.Compile(`^[a-zA-Z]+[a-zA-Z0-9._]+`)
	if err != nil {
		log.Println(err)
		return
	}
	NameExp, err = regexp.Compile(`^[a-zA-Z][a-zA-Z\s\-\']*$`)
	if err != nil {
		log.Println(err)
		return
	}
}

func ValidName(s string) bool {
	if len(s) < NAME_MIN_LEN || len(s) > NAME_MAX_LEN {
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

func ValidChatMessage(content string) bool {
	return len(content) >= CHAT_MSG_MIN_LEN && len(content) <= CHAT_MSG_MAX_LEN
}

func ValidGroupName(name string) bool {
	return len(name) >= GROUP_NAME_MIN_LEN && len(name) <= GROUP_NAME_MAX_LEN
}

func ValidGroupDescription(desc string) bool {
	return len(desc) <= GROUP_DESC_MAX_LEN
}

func ValidProfileAbout(about string) bool {
	return len(about) <= ABOUT_MAX_LEN
}

func IsAllowedImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}
