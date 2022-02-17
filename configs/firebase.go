package configs

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func SetupFirebase() *auth.Client {
	keyFilePath, err := filepath.Abs("./authentication-84a31-firebase-adminsdk-rk1wv-c1cda5c589.json")
	if err != nil {
		panic(err)
	}

	opt := option.WithCredentialsFile(keyFilePath)
	// firebase sdk inilization
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		panic(err)
	}

	// Firebase authentication
	auth, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}
	return auth
}
