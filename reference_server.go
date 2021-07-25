package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/philips-software/go-hsdp-api/console"

	"github.com/cloudfoundry-community/gautocloud"
	"github.com/dgrijalva/jwt-go"
	"github.com/philips-software/gautocloud-connectors/hsdp"

	"github.com/philips-labs/terraform-backend-hsdp/backend"
	"github.com/philips-labs/terraform-backend-hsdp/backend/store/s3"
)

func main() {
	// Config
	viper.SetEnvPrefix("tfstate")
	viper.SetDefault("key", "")
	viper.SetDefault("regions", "us-east,eu-west")
	viper.SetDefault("allow_list", "")
	viper.AutomaticEnv()

	encryptionKey := viper.GetString("key")
	hsdpRegions := strings.Split(viper.GetString("regions"), ",")
	allowList := viper.GetString("allow_list")

	if encryptionKey == "" {
		log.Printf("encryption key cannot be blank\n")
		return
	}

	// S3 bucket
	var svc *hsdp.S3MinioClient
	err := gautocloud.Inject(&svc)
	if err != nil {
		log.Printf("gautocloud: %v\n", err)
		return
	}

	// create a store
	store := s3.NewStore(&s3.Options{
		Client: svc.Client,
		Bucket: svc.Bucket,
	})

	// create a backend
	tfbackend := backend.NewBackend(store, &backend.Options{
		EncryptionKey: []byte(encryptionKey),
		Logger: func(level, message string, err error) {
			if err != nil {
				log.Printf("%s: %s - %v", level, message, err)
			} else {
				log.Printf("%s: %s", level, message)
			}
		},
		GetMetadataFunc: func(state map[string]interface{}) map[string]interface{} {
			// fmt.Println(state)
			return map[string]interface{}{
				"test": "metadata",
			}
		},
		GetRefFunc: refFunc(hsdpRegions, allowList),
	})
	if err := tfbackend.Init(); err != nil {
		log.Fatal(err)
	}

	// state
	http.HandleFunc("/states", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tfbackend.HandleListStates(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// version
	http.HandleFunc("/versions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tfbackend.HandleListVersions(w, r)
		case http.MethodDelete:
			tfbackend.HandleKeepVersions(w, r)
		case "RETRIEVE":
			tfbackend.HandleRetrieveVersion(w, r)
		case http.MethodPut:
			tfbackend.HandleRestoreVersion(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// add handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "LOCK":
			tfbackend.HandleLockState(w, r)
		case "UNLOCK":
			tfbackend.HandleUnlockState(w, r)
		case http.MethodGet:
			tfbackend.HandleGetState(w, r)
		case http.MethodPost:
			tfbackend.HandleUpdateState(w, r)
		case http.MethodDelete:
			tfbackend.HandleDeleteState(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func refFunc(regions []string, allowList string) func(*http.Request) (string, error) {
	var allowed []string
	clients := make(map[string]*console.Client, len(regions))

	for _, region := range regions {
		client, err := console.NewClient(nil, &console.Config{
			Region: region,
		})
		if err == nil {
			clients[region] = client
		}
	}
	if allowList != "" {
		allowed = strings.Split(allowList, ",")
	}
	return func(r *http.Request) (string, error) {
		// Authenticate
		username, password, ok := r.BasicAuth()
		if !ok {
			return "", fmt.Errorf("missing authentication")
		}
		if len(allowed) > 0 {
			approved := false
			for _, u := range allowed {
				if username == u {
					approved = true
					break
				}
			}
			if !approved {
				return "", fmt.Errorf("not authorized to use this backend")
			}
		}
		checkRegion := r.URL.Query().Get("region")
		authenticated := false

		var client *console.Client
		for region, rc := range clients {
			if checkRegion != "" && region != checkRegion {
				continue
			}
			c, err := rc.WithLogin(username, password)
			if err == nil && c != nil {
				client = c
				authenticated = true
				break
			}
			if c != nil {
				c.Close()
			}
		}
		if !authenticated || client == nil {
			return "", fmt.Errorf("authorization failed")
		}
		defer client.Close()
		token, _ := jwt.Parse(client.IDToken(), func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["sub"] == "" {
			return "", fmt.Errorf("invalid claims")
		}
		userUUID := claims["sub"].(string)
		path := r.URL.Path
		if path == "/versions" || path == "/states" {
			path = "/"
		}
		queryRef := r.URL.Query().Get("ref")
		if queryRef != "" {
			path = queryRef
		}
		return filepath.Join(userUUID, path), nil
	}
}
