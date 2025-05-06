package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Registry = prometheus.NewRegistry()

	Up = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Shows if the slashing observer program is up",
	})

	CountOfSlashableVaults = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "count_of_slashable_vault",
		Help: "Shows number of slashable vaults",
	})

	CountOfObservers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "count_of_observers",
		Help: "Shows number of currently active observers",
	})

	SlashingEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count_of_slashing_event",
			Help: "Count of slashing events by vault and slasher",
		},
		[]string{"vault", "slasher"},
	)
)

func Init() {
	Registry.MustRegister(Up)
	Registry.MustRegister(CountOfSlashableVaults)
	Registry.MustRegister(CountOfObservers)
	Registry.MustRegister(SlashingEventCounter)
}
