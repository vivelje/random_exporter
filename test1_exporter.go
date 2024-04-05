package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func simulateClusterState(clusterGaugeVec *prometheus.GaugeVec) {
	// Тут мы рандомизируем состояние кластеров
	for {
		for i := 1; i <= 3; i++ {
			state := 1
			if rand.Intn(100) < 5 { // В 5% случаев кластер "выключается"
				state = 0
			}
			clusterGaugeVec.With(prometheus.Labels{"cluster_id": string('0' + rune(i))}).Set(float64(state))
		}
		time.Sleep(60 * time.Second) // Обновляем значение каждую минуту
	}
}

func simulateNginxStatus(nginxStatusVec *prometheus.CounterVec) {
	// Рандомизируем статусы nginx
	for {
		statusCodes := []string{"200", "201", "400", "500"}
		counts := []int{rand.Intn(100) + 1, rand.Intn(100) + 1, rand.Intn(10), rand.Intn(5)}

		for i, code := range statusCodes {
			nginxStatusVec.With(prometheus.Labels{"status_code": code}).Add(float64(counts[i]))
		}

		time.Sleep(60 * time.Second) // Обновляем значение каждую минуту
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	clusterGaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vd_cluster_state",
			Help: "Shows whether a cluster is up (1) or down (0)",
		},
		[]string{"cluster_id"},
	)

	nginxStatusVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vd_nginx_status",
			Help: "Counts nginx responses by status code",
		},
		[]string{"status_code"},
	)

	// Регистрируем наши метрики
	prometheus.MustRegister(clusterGaugeVec)
	prometheus.MustRegister(nginxStatusVec)

	// Симуляция состояния кластера и статусов nginx в отдельных горутинах
	go simulateClusterState(clusterGaugeVec)
	go simulateNginxStatus(nginxStatusVec)

	// Запуск HTTP сервера
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
