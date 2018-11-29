# Instrumenation Library

This is the Go client library for Instrumenting go application code. This library is a wrapper over Prometheus Go Client Library, the wrapper only contains code to add metrics.
Note: The wrapper can be extended in future to push metrics data to Prometheus server. 

## Getting Started

Import this library in your code to start instrumenting your go code.
The instrumentation will start when the library is imported in your application. The library uses Prometheus go_client, the package has an init method that starts instrumenting you application code.
To host the metrics on an HTTP server in your code, you need to invoke the StartListening method, with a boolean value and logger, the boolean value indicates whether to start the httpserver in your code.

### Prerequisites

Please include the required Prometheus packages and subpackages in the glide files. Please refer to the common lib glide file for the package list.
*[glide.yaml](https://github.com/ContinuumLLC/platform-common-lib/blob/master/src/glide.yaml)
*[glide.lock](https://github.com/ContinuumLLC/platform-common-lib/blob/master/src/glide.lock)

### Code Sample
Please refer to the [main.go](https://github.com/ContinuumLLC/platform-common-lib/tree/master/src/testApps/Instrumentation/main.go) in the testApps/instrumentation package in the platform-common-lib project for sample code.


## Built With

* [Prometheus Go client library](https://github.com/prometheus/client_golang)
