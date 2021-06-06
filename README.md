# prompush
Pushing data to Prometheus remote writes from Kafka 

I can probably also describe this project as shoving Kafka in between  Prometheus and a remote storage system.

## Getting Started
* Start the docker-compose of the test services
    ```sh
    docker-compose -f test-services.yml up -d
    ```
* Build bytehoops, the Kafka proxy that will a Prometheus server will `remote_write` to
    ```shell
    cd bytehoops
    go build
    ```
* Run two instances of bytehoops, a producer and a consumer in different terminal windows
  ```shell
  ./bytehoops producer
  ```
  ```shell
  ./bytehoops consumer
  ```
* Start [avalance](https://github.com/open-fresh/avalanche) a load testing tool for Prometheus remote write
    ```shell
     docker run quay.io/freshtracks.io/avalanche --remote-url http://host.docker.internal:8080/v1/write
    ```
* See the following command to find that 5000+ timeseries were populated [cortex stats](http://localhost:9009/distributor/all_user_stats)