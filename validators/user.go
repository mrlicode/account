package validators

import "regexp"

var (
	nicknamePattern = regexp.MustCompile(`^.{2,128}$`)
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9]{6,128}$`)
	passwordPattern = regexp.MustCompile(`^[a-zA-Z0-9]{6,128}$`)
)

func ValidateNickname(value string) (pass bool, msg string) {
	if !nicknamePattern.MatchString(value) {
		return false, "仅支持大小写字母、数字，长度 2- 128 位"
	}
	return true, ""
}

func ValidateUsername(value string) (pass bool, msg string) {
	if !nicknamePattern.MatchString(value) {
		return false, "仅支持大小写字母、数字，长度 2- 128 位"
	}
	return true, ""
}

func ValidatePassword(value string) (pass bool, msg string) {
	if !nicknamePattern.MatchString(value) {
		return false, "仅支持大小写字母、数字，长度 2- 128 位"
	}
	return true, ""
}