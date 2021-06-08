# prompush
Pushing data to Prometheus remote writes from Kafka 

### Goal
To demonstrate how one can use Kafka as an intermediate buffer between Prometheus and a remote storage backend.

### Components Used
* [Bytehoops](https://github.com/tuxiedev/bytehoops) - An application that transfers bytes received on an HTTP endpoint and forwards it to Kafka. Conversely, it can also be run to consumer from Kafka and write the message to an API endpoint. In these operations, the HTTP Request Headers and Kafka Message Headers are retained
* [Apache Kafka](https://kafka.apache.org) - A high throughput messaging queue
* [Cortex](https://cortexmetrics.io) - One of the most popular remote storage and querying backends for Prometheus.
* [Avalanche](https://github.com/open-fresh/avalanche) - A load testing tool to analyse loads on a Prometheus server and remote storage backends.

### Running the experiment
* Start the docker-compose of the test services
    ```sh
    docker-compose up -d
    ```

* Start [avalance](https://github.com/open-fresh/avalanche) a load testing tool for Prometheus remote write
    ```shell
     docker run quay.io/freshtracks.io/avalanche --remote-url http://host.docker.internal:8080/v1/write
    ```
* See the following command to find that 5000+ timeseries were populated [cortex stats](http://localhost:9009/distributor/all_user_stats)
