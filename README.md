# iptools (ipt)

A command-line tool for IP address operations and CIDR range manipulations.

## Requirements

- Go 1.24.3 or higher

## Installation

```bash
go install github.com/db757/iptools@latest
```

Or build from source:

```bash
make build
```

The binary will be created in the `dist` directory.

## Usage

### Check if IP is in Range

```bash
ipt inrange [ip] [ranges]
```

### Get CIDR Range Boundaries

```bash
ipt cidrange [cidr]
```

### Get Next IP

```bash
ipt next [ip]
```

### Get Previous IP

```bash
ipt prev [ip]
```

### Get N IPs from CIDR Range

```bash
ipt getn [cidr] [count] [--offset|-o offset] [--tail|-t]
```

Options:

- `--offset, -o`: Number of IPs to skip before starting to return results
- `--tail, -t`: Count backwards from the end of the range
- `--short, -s`: Short output format (global flag)

## Development

### Available Make Commands

- `make build`: Build the project (includes tidy, clean, fmt, vet, test)
- `make test`: Run tests
- `make fmt`: Format code
- `make vet`: Run go vet
- `make lint`: Run linter
- `make clean`: Clean build artifacts
- `make targets`: Build for multiple platforms (linux-amd64, darwin-arm64)
- `make govulncheck`: Run vulnerability checks
- `make upgrade`: Upgrade dependencies
- `make run`: Run the application (includes vet)

## Dependencies

- github.com/urfave/cli/v3: CLI framework
- go4.org/netipx: IP address manipulation
- github.com/stretchr/testify: Testing framework
