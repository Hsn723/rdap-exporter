# rdap-exporter

RDAP Exporter is a Prometheus exporter for domain RDAP information. Currently it provides metrics for domain status and events such as registration and expiration dates.

## Motivation

Domain status information can be used to monitor indicators of compromise (IOC) such as the unintended transition from `clientTransferProhibited` to `pendingTransfer` hinting at an ongoing domain takeover attempt.

Domain events are useful for monitoring expiry dates for potential domain renewal issues.

## Usage

```sh
Usage:
  rdap-exporter [flags]

Flags:
  -c, --config string   path to configuration file (default "/etc/rdap-exporter/config.toml")
  -h, --help            help for rdap-exporter
```

## Configuration

```toml
domains = ["example.com", "example.net"]
check_interval = 100 # default: 60
timeout = 100        # default: 30
listen_port = 9999   # default: 9099
```

## Metrics

```sh
# HELP rdap_domain_event Dates pertaining to the domain as a unix timestamp.
# TYPE rdap_domain_event gauge
rdap_domain_event{domain="example.com",event="expiration"} 1.7235216e+09
rdap_domain_event{domain="example.com",event="last_changed"} 1.691996498e+09
rdap_domain_event{domain="example.com",event="last_update_of_rdap_database"} 1.698728968e+09
rdap_domain_event{domain="example.com",event="registration"} 8.083728e+08
# HELP rdap_domain_status Domain status codes.
# TYPE rdap_domain_status gauge
rdap_domain_status{domain="example.com",status="client_delete_prohibited"} 1
rdap_domain_status{domain="example.com",status="client_transfer_prohibited"} 1
rdap_domain_status{domain="example.com",status="client_update_prohibited"} 1
```
