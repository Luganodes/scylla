# Scylla: Symbiotic Slashing Observer
Usage: scylla <command> [flags]

### Description:

Scylla is a tool designed to monitor the Symbiotic protocol for potential slashing events. It provides insights into the health and status of Symbiotic vaults, allowing operators to react swiftly to any anomalies.

### Installation
Run the following to install scylla on your system:

```bash
curl https://raw.githubusercontent.com/Luganodes/scylla/main/install.sh | sudo bash
``` 
## Global Flags

| Flag Name           | Shorthand | Default Value                 | Description                                                                      |
| :------------------ | :-------- | :---------------------------- | :------------------------------------------------------------------------------- |
| `--help`            | `-h`      |                               | Show context-sensitive help. Displays available commands and flags.              |
| `--version`         |           |                               | Show the current version of the Scylla application.                             |
| `--ethereum.rpc`    |           | `http://127.0.0.1:8545`       | The URL of the Ethereum RPC endpoint that Scylla will connect to for on-chain data. |
| `--ethereum.ws`     |           | `http://127.0.0.1:8546`       | The URL of the Ethereum WebSocket endpoint that Scylla will subscribe to for real-time event updates. |
| `--prometheus.host` |           | `0.0.0.0`                     | The network interface address on which the Prometheus metrics server will listen for incoming scrape requests. `0.0.0.0` means it will listen on all available network interfaces. |
| `--prometheus.port` |           | `9090`                        | The TCP port number on which the Prometheus metrics server will be accessible.      |
| `--symbiotic.api`   |           | `https://app.symbiotic.fi/api/v2/vaults` | The base URL of the Symbiotic API endpoint used to fetch information about Symbiotic vaults. |
| `--config.file`     |           | `""`                          | **[OPTIONAL]** Specifies the path to a TOML configuration file. Settings defined in this file will take precedence over any command-line flags. |

### Commands:

### start
Usage: `scylla start [flags]`

#### Description:

The start command initiates the Scylla slashing observer. It begins fetching data from the specified Ethereum and Symbiotic endpoints and starts serving Prometheus metrics. This command is essential for running Scylla in the background and continuously monitoring the Symbiotic protocol.


```
scylla start --ethereum.rpc="<YOUR_RPC_ENDPOINT>" --ethereum.ws="<YOUR_WS_ENDPOINT>" --prometheus.port=9100
``` 

### How it works:

When you run scylla start, the application will:

- **Connect to Ethereum**: Establish a connection to the Ethereum network using the provided RPC and WebSocket URLs. This allows Scylla to listen for relevant on-chain events.
- **Fetch Symbiotic Data**: Periodically fetch data about Symbiotic vaults from the specified Symbiotic API endpoint.
- **Expose Prometheus Metrics**: Start a Prometheus HTTP server on the configured host and port. This server will expose various metrics related to the Symbiotic protocol and the status of its vaults. These metrics can then be scraped by a Prometheus server for monitoring and alerting.
### Metrics Served (when start is running):

When the scylla start command is executed, the application exposes the following metrics on the Prometheus endpoint (configured by --prometheus.host and --prometheus.port). These metrics provide insights into the Symbiotic protocol and potential slashing events:

- `count_of_observers` **(gauge)**: Shows the current number of active observer instances monitoring the Symbiotic protocol.
- `count_of_slashable_vault`** (gauge)**: Indicates the current number of Symbiotic vaults that are considered eligible for slashing based on the configured rules.
- `count_of_slashing_event` **(counter)**: Represents the total number of slashing events recorded by Scylla from execution. This metric is labeled with the Ethereum address of the slasher (the entity that triggered the slash) and the vault address that was slashed. The value increments each time a slashing event is detected for a specific slasher-vault pair.
- `up` **(gauge)**: A binary metric indicating the operational status of the Scylla observer. A value of 1 signifies that the Scylla process is running and healthy, while 0 indicates that it is down or experiencing issues.


### Configuration File:

Scylla supports configuration via a TOML file specified using the --config.file flag. This allows you to define all the global flags in a structured file, making it easier to manage configurations, especially for more complex setups.

Example config.toml:

```toml
[symbiotic]
api_url = "https://app.symbiotic.fi/api/v2/vaults"

[ethereum]
rpc_url = "https://mainnet.infura.io/v3/<YOUR_INFURA_KEY>"
ws_url  = "wss://mainnet.infura.io/ws/v3/<YOUR_INFURA_KEY>"

[prometheus]
host = "127.0.0.1"
port = 9090
```


```
scylla start --config.file="./config.toml"
```

