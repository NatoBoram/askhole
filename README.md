# Askhole

[![Go](https://github.com/NatoBoram/askhole/actions/workflows/go.yaml/badge.svg)](https://github.com/NatoBoram/askhole/actions/workflows/go.yaml)
[![Docker](https://github.com/NatoBoram/askhole/actions/workflows/docker.yaml/badge.svg)](https://github.com/NatoBoram/askhole/actions/workflows/docker.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/NatoBoram/askhole)](https://goreportcard.com/report/github.com/NatoBoram/askhole)

A Caddy "ask" endpoint for Kubo.

> ##### [`on_demand_tls`](https://caddyserver.com/docs/caddyfile/options#on-demand-tls)
>
> Configures [On-Demand TLS](https://caddyserver.com/docs/automatic-https#on-demand-tls) where it is enabled, but does not enable it (to enable it, use the [`on_demand` subdirective of the `tls` directive](https://caddyserver.com/docs/caddyfile/directives/tls#syntax)). Required for use in production environments, to prevent abuse.
>
> - **ask** will cause Caddy to make an HTTP request to the given URL, asking whether a domain is allowed to have a certificate issued.
>
>   The request has a query string of `?domain=` containing the value of the domain name.
>
>   If the endpoint returns a `2xx` status code, Caddy will be authorized to obtain a certificate for that name. Any other status code will result in cancelling issuance of the certificate and erroring the TLS handshake.

> [!TIP]
>
> The ask endpoint should return _as fast as possible_, in a few milliseconds, ideally. Typically, your endpoint should do a constant-time lookup in an database with an index by domain name; avoid loops. Avoid making DNS queries or other network requests.

> - **permission** allows custom modules to be used to determine whether a certificate should be issued for a particular name. The module must implement the [`caddytls.OnDemandPermission` interface](https://pkg.go.dev/github.com/caddyserver/caddy/v2/modules/caddytls#OnDemandPermission). An `http` permission module is included, which is what the `ask` option uses, and remains as a shortcut for backwards compatibility.
>
> - ⚠️ **interval** and **burst** rate limiting options were available, but are NOT recommended. Remove them from your config if you still have them.
>
> ```caddy
> {
> 	on_demand_tls {
> 		ask http://localhost:9123/ask
> 	}
> }
>
> https:// {
> 	tls {
> 		on_demand
> 	}
> }
> ```
