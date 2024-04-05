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

		for _, code := range statusCodes {
			var count int
			switch code {
			case "200", "201":
				count = rand.Intn(100) + 1 // Значения в диапазоне от 1 до 100
			case "400":
				count = rand.Intn(10) // Значения в диапазоне от 0 до 10
			case "500":
				count = rand.Intn(2) + 1 // Значения либо 1, либо 2
			}
			nginxStatusVec.With(prometheus.Labels{"status_code": code}).Add(float64(count))
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
