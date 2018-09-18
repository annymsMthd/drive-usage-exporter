package main

import (
	"net/http"
	"strings"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	used = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_drive_stats_used_bytes",
		Help: "The amount of bytes used from the disk",
	},
		[]string{"dev", "directory"})

	capacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_drive_stats_capacity_bytes",
		Help: "The amount of bytes on the disk",
	},
		[]string{"dev", "directory"})
)

func main() {
	prometheus.MustRegister(used)
	prometheus.MustRegister(capacity)

	go updateStats()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func updateStats() {
	for {
		fs := sigar.FileSystemList{}
		err := fs.Get()
		if err != nil {
			panic(err)
		}

		for _, fs := range fs.List {
			d := fs.DirName

			usage := sigar.FileSystemUsage{}

			err = usage.Get(d)
			if err != nil {
				panic(err)
			}

			if !strings.HasPrefix(fs.DevName, "/dev/sd") && !strings.HasPrefix(fs.DevName, "/dev/mapper") {
				continue
			}

			used.WithLabelValues(fs.DevName, d).Set(float64(usage.Used))
			capacity.WithLabelValues(fs.DevName, d).Set(float64(usage.Total))
		}

		time.Sleep(5 * time.Second)
	}
}
