package utils

import "errors"

var (
	ErrValidatePlatform = errors.New("error validate platform")
)

func ValidatePlatform(platform string) error {
	if platform != "pc-desktop" && platform != "web" && platform != "ios-mobile" && platform != "android-mobile" {
		return ErrValidatePlatform
	}
	return nil
}