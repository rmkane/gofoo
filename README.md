# Go Foo

A Cobra CLI template that handles loading Viper config.

## Table of Contents

- [Go Foo](#go-foo)
  - [Table of Contents](#table-of-contents)
  - [Init](#init)
  - [Commands](#commands)
    - [`init`](#init-1)
    - [`show`](#show)
  - [Logging](#logging)
  - [Configuration](#configuration)
    - [Example Configuration](#example-configuration)
  - [Usage](#usage)
  - [License](#license)

## Init

The `init` command initializes the configuration by checking for a configuration file in the following order of precedence:

1. If a configuration file is specified using the `--config` or `-c` flag, use that.
2. If a configuration file exists next to the binary, use that.
3. If a configuration file exists in the current working directory, use that.
4. If a configuration file exists in the home directory, use that.

If no configuration file is found, the default configuration is loaded.

For checks 1-4, the configuration file can have the following extensions:

- `.yaml` or `.yml`
- `.json`
- `.toml`

## Commands

### `init`

Initializes the configuration file.

```sh
gofoo init --format yaml --force
```

- `--format` or `-f`: Specifies the format of the configuration file. Supported formats are json, yaml, and toml. Default is yaml.
- `--force`: Overwrites the configuration file if it exists

### `show`

Displays the current configuration.

```sh
gofoo show
```

## Logging

The application uses `slog` for logging. The log level and format can be configured through the configuration file.

## Configuration

The configuration file supports the following keys:

- `logging.dir`: Directory where log files are stored.
- `logging.format`: Format of the log files. Supported formats are `json` and `text`.
- `logging.level`: Log level. Supported levels are `DEBUG`, `INFO`, `WARN`, and `ERROR`.

### Example Configuration

```yml
logging:
  dir: /var/log/gofoo
  format: json
  level: debug
```

## Usage

To run the application, use the following command:

```sh
go run cmd/gofoo/main.go
```

To specify a configuration file, use the `--config` or `-c` flag:

```sh
go run cmd/gofoo/main.go --config /path/to/config.yml
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
