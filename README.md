# loglocation2source

A service transforms VictoriaMetrics log metric labels app_version and location into a GitHub Source Code link. 

It is used as [data link](https://grafana.com/docs/grafana/latest/visualizations/panels-visualizations/configure-data-links/) in VictoriaMetrics Grafana dashboard to quickly navigate from logs panel to the source code location on GitHub.

Try it out:
- [Tag](https://loglocation2source.makasim.com/?app_version=vmagent-20251104-105304-tags-v1.129.1-0-g5e98e0cff5&location=VictoriaMetrics/lib/promscrape/scrapework.go:394)
- [Dirty Tag](https://loglocation2source.makasim.com/?app_version=vmagent-20251104-105304-tags-v1.129.1-0-g5e98e0cff5-dirty-123abc&location=VictoriaMetrics/lib/promscrape/scrapework.go:394)
