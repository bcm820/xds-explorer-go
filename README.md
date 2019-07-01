# xDS Explorer

xDS Explorer provides a RESTful interface for viewing the current state of resources that are discovered by an [Envoy](https://www.envoyproxy.io/) management server and exposed via its [Aggregated Discovery Service (ADS)](https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/v2_overview#aggregated-discovery-service).

Since ADS is only available via gRPC, it is unavailable client-side via web browsers for two reasons:
1. A gRPC connection must be established with the Envoy management server with a [DiscoveryRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/discovery.proto#envoy-api-msg-discoveryrequest).
2. The [DiscoveryResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/discovery.proto#envoy-api-msg-discoveryresponse) returned in a gRPC connection stream is a protobuf message that must be marshaled into JSON.

So this service simply establishes the gRPC connection, creates the DiscoveryRequest with the given inputs, and returns the DiscoveryResponse as JSON.

![demo](./screencap.gif "demo")

## Resource Types

ADS exposes [Envoy resource types](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol#type-urls), all of which are available via xDS Explorer:
* Clusters (CDS)
* Cluster Load Assignments, e.g. Endpoints (EDS)
* Routes (RDS)
* Listeners (LDS)
* Secrets (SDS)

## Basic Usage

* `docker pull bcmendoza/xds-explorer` to get the Docker image or `make build` to create `bcmendoza/xds-explorer:latest`.
* Add the image to your docker-compose setup with the Envoy management server. Expose port 3001 to your local machine. Set environment variables `XDS_HOST` and `XDS_PORT`.

```yaml
  xds-explorer:
    image: bcmendoza/xds-explorer:latest
    environment:
      - XDS_HOST=gm-control
      - XDS_PORT=50000
    ports:
      - "3001:3001"
```

* Visit `localhost:3001` in your browser.
* Select a `ResourceType` and fill out fields for `Node`, `Zone`, `Clusters`, and optionally `ResourceNames` (a comma-separated list). For example:

```
ResourceType: "ClusterLoadAssignment"

Node: "default-node"

Zone: "default-zone"

Cluster: "catalog"

ResourceNames: ""
```

## API

### Route: `/request`

#### Method: `POST`

Initializes or updates the current DiscoveryRequest being made to the Envoy management server.

Sample Request Body:
```json
{
  "resourceType": "ClusterLoadAssignment",
  "node": "default-node",
  "zone": "default-zone",
  "cluster": "catalog",
  "resourceNames": ["catalog"]
}
```

Response:
```json
{"request updated": true}
```

---

### Route: `/listen`

#### Method: `GET`

Retrieves the current DiscoveryResponse being made to the Envoy management server (or an empty array if no request has been made).

Response:
```json
[
  {
    "cluster_name": "catalog",
    "endpoints": [
      {
        "lb_endpoints": [
          {
            "HostIdentifier": {
              "Endpoint": {
                "address": {
                  "Address": {
                    "SocketAddress": {
                      "address": "192.168.16.8",
                      "PortSpecifier": {
                        "PortValue": 8080
                      }
                    }
                  }
                }
              }
            },
            "health_status": 1,
            "metadata": {}
          }
        ]
      }
    ]
  }
]
```
