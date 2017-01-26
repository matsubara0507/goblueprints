package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

import gomniauthtest "github.com/stretchr/gomniauth/test"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、",
			"AuthAvatar.GetAvatarURL は ErrNoAvatarURL を返すべきです")
	}

	testUrl := "http://url-to-avatar"
	testUser = &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return(testUrl, ErrNoAvatarURL)
	testChatUser.User = testUser
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("値が存在する場合、",
			"AuthAvatar.GetAvatarURL はエラーを返すべきではありません: ", err)
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL は正しいURLを返すべきです: ", url)
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL はエラーを返すべきではありません: ", err)
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL が %s という誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL はエラーを返すべきではありません: ", err)
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL が %s という誤った値を返しました", url)
	}
}
