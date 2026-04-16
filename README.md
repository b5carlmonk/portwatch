# portwatch

Lightweight CLI to monitor and alert on open ports and service changes on a host.

## Installation

```bash
go install github.com/youruser/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start monitoring open ports on the local host and get alerted when changes are detected:

```bash
portwatch watch
```

Specify a custom scan interval and target host:

```bash
portwatch watch --host 192.168.1.10 --interval 30s
```

Run a one-time snapshot of currently open ports:

```bash
portwatch scan
```

Example output when a change is detected:

```
[ALERT] New port opened: 8080/tcp (2024-06-10 14:32:01)
[ALERT] Port closed:     3306/tcp (2024-06-10 14:35:44)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `localhost` | Target host to monitor |
| `--interval` | `60s` | Polling interval |
| `--alert` | `stdout` | Alert output (`stdout`, `file`, `webhook`) |
| `--output` | `text` | Output format (`text`, `json`) |

## License

MIT © youruser