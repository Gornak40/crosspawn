package main

import (
	"time"

	"github.com/Gornak40/crosspawn/config"
	"github.com/Gornak40/crosspawn/controller"
	"github.com/Gornak40/crosspawn/models"
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	db, err := models.ConnectDatabase(&cfg.DBConfig)
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to database")
	}

	ej := ejudge.NewEjClient(&cfg.EjConfig)

	s := controller.NewServer(db, ej, &cfg.ServerConfig)

	for {
		if err := s.Poll(); err != nil {
			logrus.WithError(err).Error("failed to poll ejudge")
		}
		time.Sleep(time.Second * time.Duration(cfg.PollDelaySeconds))
	}
}
