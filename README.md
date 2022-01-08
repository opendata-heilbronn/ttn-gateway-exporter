# ttn-gateway-exporter
Exports The Things Network gateway status as Prometheus metrics

## Quick Start
1. Create a file with the following contents and name it targets.yaml. This assumes you're running in the eu1 community cluster. 
```yaml
targets:
  - gateway_id: my-ttn-gateway # Your TTN Gateway ID
    api_key: NNSXS.[...redacted...] # Your TTN API Key that has access to this Gateway
```

2. Run `docker run -v $PWD/targets.yaml:/etc/ttn-exporter/targets.yaml -p 8080:8080 ghcr.io/opendata-heilbronn/ttn-gateway-exporter:v0.0.1`
3. Access the exported metrics at `http://localhost:8080/metrics`
4. Point your Prometheus instance at the exporter
```yaml
scrape_configs:
  - job_name: ttn-gateway
    static_configs:
      - targets:
          - my-ttn-exporter-host.example.com:8080
```

## Docker-Compose

See the example `docker-compose.yaml`

## Usage
`ttn-gateway-exporter [--address ip:port] [--target-config-path /path/to/target/config.yaml]`

The `--address` parameter changes the IP and port where the TTN gateway exporter binds and exposed the metrics. The 
default value `:8080` will bind to port 8080 on all interfaces. You can specify the IP address of an interface to only
listen to incoming requests on this IP. For example, to listen only on localhost, you can specify `127.0.0.1:8080`.
Changing the port is as easy as replacing the value behind the colon.

The `--target-config-path` parameter defaults to `/etc/ttn-exporter/targets.yaml` and expects a YAML file with the 
following structure:

```yaml
default_base_url: https://eu1.cloud.thethings.network # The API host that should apply to all targets if nothing else is specified. Defaults to the community eu1 cluster if not specified
targets:
  - gateway_id: my-ttn-gateway # Your TTN Gateway ID
    api_key: NNSXS.[...redacted...] # Your TTN API Key that has access to this Gateway
    base_url: https://eu1.cloud.thethings.network # The API host that this Gateway is registered on. Overrides the default_base_url setting for this Gateway
  - gateway_id: my-north-america-ttn-gateway
    api_key: NNSXS.[...redacted...]
    base_url: https://nam1.cloud.thethings.network # This overrides the default_base_url to the community North America cluster
  - gateway_id: my-second-ttn-gateway
    api_key: NNSXS.[...redacted...]
    # base_url intentionally not specified, to use th default_base_url from above
```

The Docker image uses the same defaults. That means, if you want mount your config file into the Docker container, mount it to `/etc/ttn-exporter/targets.yaml`.
