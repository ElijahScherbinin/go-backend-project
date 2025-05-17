package prometheus

import (
	"user-service/internal/user"

	"github.com/prometheus/client_golang/prometheus"
)

type UserPrometheusImpl struct {
	new    prometheus.Counter
	delete prometheus.Counter
}

func (userPrometheus *UserPrometheusImpl) New() {
	userPrometheus.new.Add(1)
}

func (userPrometheus *UserPrometheusImpl) Delete() {
	userPrometheus.delete.Add(1)
}

func NewUserPrometheus() user.UserPrometheus {
	userPrometheus := UserPrometheusImpl{
		new: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "user_service",
				Name:      "new_user",
				Help:      "This is new user counter",
			},
		),
		delete: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "user_service",
				Name:      "delete_user",
				Help:      "This is delete user counter",
			},
		),
	}
	prometheus.MustRegister(userPrometheus.new)
	prometheus.MustRegister(userPrometheus.delete)
	return &userPrometheus
}
