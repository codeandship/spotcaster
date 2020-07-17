package metrics

import "github.com/prometheus/client_golang/prometheus"

type Prometheus struct {
	episodes  prometheus.Gauge
	starts    prometheus.Gauge
	streams   prometheus.Gauge
	listeners prometheus.Gauge
	followers prometheus.Gauge
}

func NewPrometheus() *Prometheus {
	p := Prometheus{
		episodes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "spotcaster",
				Subsystem: "",
				Name:      "episodes_total",
				Help:      "total count of episodes published",
			},
		),
		starts: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "spotcaster",
				Subsystem: "",
				Name:      "starts_total",
				Help:      "Measured when a Spotify user listens to 0 seconds or more of any episode in your catalog.",
			},
		),
		streams: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "spotcaster",
				Subsystem: "",
				Name:      "streams_total",
				Help:      "Measured when a Spotify user listens to 60 seconds or more of any episode in your catalog.",
			},
		),
		listeners: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "spotcaster",
				Subsystem: "",
				Name:      "listeners_total",
				Help:      "Measures the number of unique Spotify users who started an episode in your catalog.",
			},
		),
		followers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "spotcaster",
				Subsystem: "",
				Name:      "followers_total",
				Help:      "Followers are listeners who hit Follow on your podcast on Spotify.",
			},
		),
	}

	prometheus.MustRegister(
		p.episodes,
		p.followers,
		p.listeners,
		p.starts,
		p.streams,
	)
	return &p
}

func (p *Prometheus) SetEpisodes(n float64) {
	p.episodes.Set(n)
}

func (p *Prometheus) SetFollowers(n float64) {
	p.followers.Set(n)
}

func (p *Prometheus) SetListeners(n float64) {
	p.listeners.Set(n)
}

func (p *Prometheus) SetStarts(n float64) {
	p.starts.Set(n)
}

func (p *Prometheus) SetStreams(n float64) {
	p.streams.Set(n)
}
