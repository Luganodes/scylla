package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)

	handler := promhttp.HandlerFor(Registry, promhttp.HandlerOpts{})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		handler.ServeHTTP(w, r)
	})

	go func() {
		fmt.Println("📡 Prometheus metrics available at:", addr+"/metrics")
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println("❌ Prometheus server failed:", err)
		}
	}()
}
