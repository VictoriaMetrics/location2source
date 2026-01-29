# loglocation2source
A service transofrms VictoriaMetrics log metric labels into a Github Source Code link. 

Could be used as data link in Grafana dashboard to quickly navigate to the source code location from log metric.

```
topk_max(
  50,
  sum(rate(vm_log_messages_total{job=~"$job",instance=~"$instance",level!="info"}[$__rate_interval])) by (job,instance,app_version,location)
)
```

```
https://loglocation2source.makasim.com/?app_version=${__field.labels.app_version}&location=${__field.labels.location}
```

Try it out:

- Clean: https://loglocation2source.makasim.com/?app_version=vmagent-20251104-105304-tags-v1.129.1-0-g5e98e0cff5&location=VictoriaMetrics/lib/promscrape/scrapework.go:394
- Dirty: https://loglocation2source.makasim.com/?app_version=vmagent-20251104-105304-tags-v1.129.1-0-g5e98e0cff5-dirty-123abc&location=VictoriaMetrics/lib/promscrape/scrapework.go:394


