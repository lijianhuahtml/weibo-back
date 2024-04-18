package utils

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	// 正则表达式来匹配邮箱格式
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(regex, email)
	if err != nil {
		return false
	}
	return match
}

func IsValidPassword(password string) bool {
	// 正则表达式来匹配邮箱格式
	regex := `^[a-zA-Z0-9]{6,30}$`
	match, err := regexp.MatchString(regex, password)
	if err != nil {
		return false
	}
	return match
}
