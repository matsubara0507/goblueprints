package main

import (
	"errors"
	"path/filepath"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	matches, err := filepath.Glob(filepath.Join("avatars", u.UniqueID()+".*"))
	if err != nil || len(matches) == 0 {
		return "", ErrNoAvatarURL
	}
	return "/" + filepath.ToSlash(matches[0]), nil
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		url, err := avatar.GetAvatarURL(u)
		if err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
