package main

import (
	"../../chapter1/trace"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"gopkg.in/yaml.v2"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type Keys struct {
	SecurityKey   string `yaml:"securityKey"`
	OauthServices map[string]struct {
		ClientId  string `yaml:"id"`
		SecretKey string `yaml:"secret"`
	} `yaml:"oauthServices"`
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}
	authCookie, err := r.Cookie("auth")
	if err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	err = t.templ.Execute(w, data)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	buf, err := ioutil.ReadFile(".keys.yaml")
	if err != nil {
		log.Fatal("keys.yaml is not found.")
	}
	keys := Keys{}
	err = yaml.Unmarshal(buf, &keys)
	if err != nil {
		log.Fatal("error Unmarshal: ", buf)
	}

	gomniauth.SetSecurityKey(keys.SecurityKey)
	gomniauth.WithProviders(
		facebook.New(keys.OauthServices["facebook"].ClientId, keys.OauthServices["facebook"].SecretKey, "http://localhost:8080/auth/callback/facebook"),
		github.New(keys.OauthServices["github"].ClientId, keys.OauthServices["github"].SecretKey, "http://localhost:8080/auth/callback/github"),
		google.New(keys.OauthServices["google"].ClientId, keys.OauthServices["google"].SecretKey, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	go r.run()
	log.Println("Webサーバーを開始します。ポート: ", *addr)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
