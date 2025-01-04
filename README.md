# soa-404

## How to run ?

Simply run `make run` and `curl "http://localhost:8070/random?seed=123"`.

## Access Grafana

Go to `http://localhost:9000` -> Explore

Change the datasources to `Mimir` to access the metrics, `Loki` to access the logs and `Tempo` for traces.

### Metrics

<figure>
<p align="center" width="100%">
  <img src="./docs/images/metrics.png" alt="init_result" style="width:100%">
  <figcaption><p align="center">Explore metrics</p></figcaption>
  </p>
</figure>

### Logs

<figure>
<p align="center" width="100%">
  <img src="./docs/images/logs.png" alt="init_result" style="width:100%">
  <figcaption><p align="center">Explore logs</p></figcaption>
  </p>
</figure>

### Traces

<figure>
<p align="center" width="100%">
  <img src="./docs/images/traces.png" alt="init_result" style="width:100%">
  <figcaption><p align="center">Explore traces</p></figcaption>
  </p>
</figure>

## Contributing

Feel free to fork or clone this repository, explore the code, and contribute by submitting pull requests. Contributions, whether they’re bug fixes, improvements, or new features, are always welcome!

## License

Distributed under the GPLv3 License. See `LICENSE.md` file for more information.
