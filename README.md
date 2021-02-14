# cisco_telnet_exporter
Prometheus Exporter for Cisco Catalyst switches using telnet.

### Tested Devices
- Cisco 3560g-48ps

## Installation
The exporter listens on port `9504` by default.

### Docker
```
docker run -d \
    --network host \
    -e CISCO_IP='192.168.1.145' \
    -e CISCO_PORT='23' \
    -e CISCO_PASS='admin' \
    chhaley/cisco_telnet_exporter
```

### Helm
```
helm repo add charlie-haley https://charlie-haley.github.io/private-charts
helm repo update
helm install cisco-telnet-exporter charlie-haley/cisco-telnet-exporter \
    --set "cisco.ip=192.168.1.145" \ 
    --set "cisco.port=23" \
    --set "cisco.password=admin" \
    -n monitoring
```

If you want to use the ServiceMonitor (which is enabled by default) you'll need to have [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) deployed to your cluster, see [values](charts/tplink-smartplug-exporter/values.yaml) to disable it if you'd like use ingress instead.

[You can find the chart repo here](https://github.com/charlie-haley/private-charts), if you'd like to contribute. 

## Metrics
Name               | Description                          | Labels
-------------------|--------------------------------------|------
cisco_telnet_temp  | The temperature of the switch        | instance
cisco_power_used   | Current power usage in watts         | instance
