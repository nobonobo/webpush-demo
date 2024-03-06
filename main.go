package main

import (
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

//go:generate sh -c "npm run build"

//go:embed all:build/*
var assets embed.FS

var config = struct {
	Subscriber      string `env:"Subscriber,required"`
	VAPIDPublicKey  string `env:"VAPIDPublicKey,required"`
	VAPIDPrivateKey string `env:"VAPIDPrivateKey,required"`
	ClientID        string `env:"ClientID,required"`
	ClientSecret    string `env:"ClientSecret,required"`
}{}

func init() {
	log.SetFlags(log.Lshortfile)
	var dotenv string
	flag.StringVar(&dotenv, "env", ".env", "load .env file")
	flag.Parse()
	if err := godotenv.Load(dotenv); err != nil {
		log.Print(err)
	}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	hasher := sha1.New()
	hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	fp, err := os.Create(filepath.Join("subscribes", sha+".json"))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer fp.Close()
	if _, err := fp.Write(b); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"error": nil,
	}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func notify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok || clientID != config.ClientID || clientSecret != config.ClientSecret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	log.Println("push:", string(b))
	if err := fs.WalkDir(os.DirFS("subscribes"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		fp, err := os.Open(filepath.Join("subscribes", path))
		if err != nil {
			return err
		}
		defer fp.Close()
		var s *webpush.Subscription
		if err := json.NewDecoder(fp).Decode(&s); err != nil {
			return err
		}
		resp, err := webpush.SendNotificationWithContext(r.Context(), b, s, &webpush.Options{
			Subscriber:      config.Subscriber,
			VAPIDPublicKey:  config.VAPIDPublicKey,
			VAPIDPrivateKey: config.VAPIDPrivateKey,
			TTL:             30,
		})
		if err != nil {
			return err
		}
		log.Println(resp)
		return nil
	}); err != nil {
		log.Println(err)
	}
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	os.MkdirAll("subscribes", 0755)
	sub, err := fs.Sub(assets, "build")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(sub)))
	http.HandleFunc("/subscribe", subscribe)
	http.HandleFunc("/notify", notify)
	log.Fatal(http.ListenAndServe(":8080", logger(http.DefaultServeMux)))
}
