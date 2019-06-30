# xDS Explorer

xDS Explorer provides RESTful access to the state of resources discovered by an Envoy management server and exposed via its Aggregated Discovery Service (ADS).

Since ADS is only available via gRPC, it is generally unavailable client-side via web browsers for two reasons:
1. A gRPC connection must be established with the Envoy management server with a DiscoveryRequest.
2. The DiscoveryResponse returned in a gRPC connection stream is a protobuf message that must be marshaled into JSON.

So this service simply establishes the gRPC connection, creates the DiscoveryRequest with the given inputs, and returns the DiscoveryResponse marshaled as JSON.

## Resource Types

ADS exposes Envoy resource types, all of which are available via xDS Explorer:
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

### Route: `/`

#### Method: `GET`

Loads a single-page application interface for interacting with the API.

### Route: `/request`

#### Method: `POST`

Initializes or updates the current DiscoveryRequest being made to the Envoy management server.

Sample Request Body:
```json
{
  "resourceType": "RouteConfiguration",
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

### Route: `/listen`

#### Method: `GET`

Retrieves the current DiscoveryResponse being made to the Envoy management server (or an empty array if no request has been made).

Response:
```json
[
  {
    "name": "catalog:8080",
    "address": {
      "Address": {
        "SocketAddress": {
          "address": "0.0.0.0",
          "PortSpecifier": {
            "PortValue": 8080
          }
        }
      }
    },
    "filter_chains": [
      {
        "filter_chain_match": {},
        "filters": [
          {
            "name": "envoy.http_connection_manager",
            "ConfigType": {
              "Config": {
                "fields": {
                  "access_log": {
                    "Kind": {
                      "ListValue": {
                        "values": [
                          {
                            "Kind": {
                              "StructValue": {
                                "fields": {
                                  "config": {
                                    "Kind": {
                                      "StructValue": {
                                        "fields": {
                                          "additional_request_headers_to_log": {
                                            "Kind": {
                                              "ListValue": {
                                                "values": [
                                                  {
                                                    "Kind": {
                                                      "StringValue": "X-TBN-DOMAIN"
                                                    }
                                                  },
                                                  {
                                                    "Kind": {
                                                      "StringValue": "X-TBN-ROUTE"
                                                    }
                                                  },
                                                  {
                                                    "Kind": {
                                                      "StringValue": "X-TBN-RULE"
                                                    }
                                                  },
                                                  {
                                                    "Kind": {
                                                      "StringValue": "X-TBN-SHARED-RULES"
                                                    }
                                                  },
                                                  {
                                                    "Kind": {
                                                      "StringValue": "X-TBN-CONSTRAINT"
                                                    }
                                                  }
                                                ]
                                              }
                                            }
                                          },
                                          "common_config": {
                                            "Kind": {
                                              "StructValue": {
                                                "fields": {
                                                  "grpc_service": {
                                                    "Kind": {
                                                      "StructValue": {
                                                        "fields": {
                                                          "envoy_grpc": {
                                                            "Kind": {
                                                              "StructValue": {
                                                                "fields": {
                                                                  "cluster_name": {
                                                                    "Kind": {
                                                                      "StringValue": "xds_cluster"
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  },
                                                  "log_name": {
                                                    "Kind": {
                                                      "StringValue": "tbn.access"
                                                    }
                                                  }
                                                }
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  },
                                  "name": {
                                    "Kind": {
                                      "StringValue": "envoy.http_grpc_access_log"
                                    }
                                  }
                                }
                              }
                            }
                          }
                        ]
                      }
                    }
                  },
                  "http_filters": {
                    "Kind": {
                      "ListValue": {
                        "values": [
                          {
                            "Kind": {
                              "StructValue": {
                                "fields": {
                                  "config": {
                                    "Kind": {
                                      "StructValue": {}
                                    }
                                  },
                                  "name": {
                                    "Kind": {
                                      "StringValue": "envoy.cors"
                                    }
                                  }
                                }
                              }
                            }
                          },
                          {
                            "Kind": {
                              "StructValue": {
                                "fields": {
                                  "config": {
                                    "Kind": {
                                      "StructValue": {
                                        "fields": {
                                          "upstream_log": {
                                            "Kind": {
                                              "ListValue": {
                                                "values": [
                                                  {
                                                    "Kind": {
                                                      "StructValue": {
                                                        "fields": {
                                                          "config": {
                                                            "Kind": {
                                                              "StructValue": {
                                                                "fields": {
                                                                  "additional_request_headers_to_log": {
                                                                    "Kind": {
                                                                      "ListValue": {
                                                                        "values": [
                                                                          {
                                                                            "Kind": {
                                                                              "StringValue": "X-TBN-DOMAIN"
                                                                            }
                                                                          },
                                                                          {
                                                                            "Kind": {
                                                                              "StringValue": "X-TBN-ROUTE"
                                                                            }
                                                                          },
                                                                          {
                                                                            "Kind": {
                                                                              "StringValue": "X-TBN-RULE"
                                                                            }
                                                                          },
                                                                          {
                                                                            "Kind": {
                                                                              "StringValue": "X-TBN-SHARED-RULES"
                                                                            }
                                                                          },
                                                                          {
                                                                            "Kind": {
                                                                              "StringValue": "X-TBN-CONSTRAINT"
                                                                            }
                                                                          }
                                                                        ]
                                                                      }
                                                                    }
                                                                  },
                                                                  "common_config": {
                                                                    "Kind": {
                                                                      "StructValue": {
                                                                        "fields": {
                                                                          "grpc_service": {
                                                                            "Kind": {
                                                                              "StructValue": {
                                                                                "fields": {
                                                                                  "envoy_grpc": {
                                                                                    "Kind": {
                                                                                      "StructValue": {
                                                                                        "fields": {
                                                                                          "cluster_name": {
                                                                                            "Kind": {
                                                                                              "StringValue": "xds_cluster"
                                                                                            }
                                                                                          }
                                                                                        }
                                                                                      }
                                                                                    }
                                                                                  }
                                                                                }
                                                                              }
                                                                            }
                                                                          },
                                                                          "log_name": {
                                                                            "Kind": {
                                                                              "StringValue": "tbn.upstream"
                                                                            }
                                                                          }
                                                                        }
                                                                      }
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          },
                                                          "name": {
                                                            "Kind": {
                                                              "StringValue": "envoy.http_grpc_access_log"
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  }
                                                ]
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  },
                                  "name": {
                                    "Kind": {
                                      "StringValue": "envoy.router"
                                    }
                                  }
                                }
                              }
                            }
                          }
                        ]
                      }
                    }
                  },
                  "rds": {
                    "Kind": {
                      "StructValue": {
                        "fields": {
                          "config_source": {
                            "Kind": {
                              "StructValue": {
                                "fields": {
                                  "api_config_source": {
                                    "Kind": {
                                      "StructValue": {
                                        "fields": {
                                          "api_type": {
                                            "Kind": {
                                              "StringValue": "GRPC"
                                            }
                                          },
                                          "grpc_services": {
                                            "Kind": {
                                              "ListValue": {
                                                "values": [
                                                  {
                                                    "Kind": {
                                                      "StructValue": {
                                                        "fields": {
                                                          "envoy_grpc": {
                                                            "Kind": {
                                                              "StructValue": {
                                                                "fields": {
                                                                  "cluster_name": {
                                                                    "Kind": {
                                                                      "StringValue": "xds_cluster"
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  }
                                                ]
                                              }
                                            }
                                          },
                                          "refresh_delay": {
                                            "Kind": {
                                              "StringValue": "30s"
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          },
                          "route_config_name": {
                            "Kind": {
                              "StringValue": "catalog:8080"
                            }
                          }
                        }
                      }
                    }
                  },
                  "stat_prefix": {
                    "Kind": {
                      "StringValue": "catalog-8080"
                    }
                  }
                }
              }
            }
          }
        ]
      }
    ],
    "listener_filters": null
  }
]
```
