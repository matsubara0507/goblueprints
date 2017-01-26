package main

import (
	"errors"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if ok {
		urlStr, ok := url.(string)
		if ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]
	if ok {
		useridStr, ok := userid.(string)
		if ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
