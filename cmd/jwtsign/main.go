package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Gornak40/crosspawn/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const (
	defaultUser     = "babayka"
	defaultDuration = 1 * time.Hour
)

func main() {
	user := flag.String("user", defaultUser, "ejudge login for JWT")
	duration := flag.Duration("duration", defaultDuration, "duration for JWT")
	flag.Parse()

	logrus.WithFields(logrus.Fields{
		"user":     *user,
		"duration": *duration,
	}).Info("generating JWT")

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	key := cfg.JWTSecret
	if key == "" {
		logrus.Fatal("JWT_SECRET env var is not set")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "crosspawn",
		"sub": *user,
		"exp": time.Now().Add(*duration).Unix(),
	})
	s, err := t.SignedString([]byte(key))
	if err != nil {
		logrus.WithError(err).Fatal("failed to sign JWT")
	}

	fmt.Println(s) //nolint:forbidigo // basic functionality
}
