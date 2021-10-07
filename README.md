# Reproduction

```sh
git clone https://github.com/open-telemetry/opentelemetry-collector-contrib.git
cd opentelemetry-collector-contrib
make otelcontribcol
./bin/otelcontribcol_linux_amd64 --config ../../github.com/tyrone-anz/export-otlp-googlecloud/collector-config.yaml
```

In a separate terminal, in tyrone-anz/export-otlp-googlecloud:
```sh
go run main.go
```
