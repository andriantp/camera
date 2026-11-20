# Camera

A lightweight Go library for controlling IP cameras using **ONVIF** and **Hikvision ISAPI**.

Simple and consistent Go API without exposing protocol-specific XML or HTTP details.

## Features

### ONVIF

- GetCapabilities
- GetConfigurations
- GetProfiles
- GetStatus
- RelativeMove

### Hikvision ISAPI

- GetInputs
- GetStreaming
- GetCapabilities
- GetDeviceInfo
- GetSystemCapabilities
- GetStreamingCapabilities
- GetPTZCapabilities
- GetEventCapabilities
- GetImageCapabilities
- GetAnalyticsCapabilities
- GetStatus
- GetAbsolute
- SetAbsolute
- SetZoom
- GetPresets
- GotoPreset
- GotoHome
- RelativeMove
- GetSnapshot
- GetContentCapabilities
- GetStorage
- ContentSearch
- ContentDownload
- SetThermometry
- GetThermometryRulesTemperatureInfo
- Reboot
- GetHTTPHosts
- GetMotionDetection

## Installation

```bash
go get github.com/andriantp/camera
```

## Quick Start

### Configuration

```go
client := camera.New(camera.Config{
    Timeout:         30 * time.Second, // HTTP timeout per request (default 30s)
    DownloadTimeout: 5 * time.Minute,  // HTTP timeout for recording downloads (default 5m)
    MoveTimeout:     10 * time.Second, // how long SetAbsolute/SetZoom wait for the camera to reach the target (default 5s)
})
```

All fields are optional. A shorter `context` deadline always takes precedence.
When `MoveTimeout` elapses, the error wraps `isapi.ErrMoveTimeout`:

```go
if errors.Is(err, isapi.ErrMoveTimeout) {
    // camera did not reach the target in time
}
```

### HTTPS

Prefix the camera host with `https://` to use HTTPS for the ISAPI driver
(a host without a scheme uses HTTP, as before):

```go
isapi.Camera{
    Credential: camera.Credential{
        Host:     "https://192.168.1.64", // or "https://192.168.1.64:8443"
        Username: "...",
        Password: "...",
    },
    Channel: 1,
}
```

Cameras usually ship a self-signed certificate, which is rejected by default.
Trust it explicitly with `Config.TLSConfig`:

```go
pem, _ := os.ReadFile("camera.crt") // exported from the camera web UI
roots := x509.NewCertPool()
roots.AppendCertsFromPEM(pem)

client := camera.New(camera.Config{
    TLSConfig: &tls.Config{RootCAs: roots},
})
```

`&tls.Config{InsecureSkipVerify: true}` also works, but disables certificate
checks and exposes credentials to man-in-the-middle attacks; use it only on a
trusted network.

The ONVIF driver always uses HTTP (the underlying library hardcodes it). An
`https://` host is accepted so one `Credential` can be shared with ISAPI: the
scheme and HTTPS port are dropped and ONVIF connects to the default HTTP port.
If ONVIF listens on a non-default port, give the ONVIF client its own
`Credential` with `"host:port"`.

### ONVIF

```go
client := camera.New(camera.Config{
    Timeout: 30 * time.Second,
})

onvifClient := onvif.New(client)

capabilities, err := onvifClient.GetCapabilities(
    context.Background(),
    onvif.Camera{
        Credential: camera.Credential{
            Host:     "...",
            Username: "...",
            Password: "...",
        },
    },
)
if err != nil {
    log.Fatal(err)
}
```

### Hikvision ISAPI

```go
client := camera.New(camera.Config{
    Timeout: 30 * time.Second,
})

isapiClient := isapi.New(client)

capabilities, err := isapiClient.GetCapabilities(
    context.Background(),
    isapi.Camera{
        Credential: camera.Credential{
            Host:     "...",
            Username: "...",
            Password: "...",
        },
        Channel: 1,
    },
)
if err != nil {
    log.Fatal(err)
}
```

## Project Structure

```text
.
├── client.go
├── driver/
│   ├── isapi/
│   └── onvif/
├── model/
├── examples/
└── README.md
```

- `driver` - Protocol implementations
- `model` - Shared public models
- `examples` - Usage examples

## Examples

```text
examples/
├── onvif/
└── isapi/
```

Run an example:

```bash
go run examples/onvif/status/main.go
```

or:

```bash
go run examples/isapi/channels/main.go INPUTS
```

## Design

```text
Application
    ↓
Driver
    ↓
HTTP / Protocol
    ↓
XML Parser
    ↓
Public Model
```

Protocol-specific implementation stays inside the driver, while applications use shared Go models.

## Supported Protocols

| Protocol | Status |
|---|---|
| ONVIF | ✅ |
| Hikvision ISAPI | ✅ |

## License

MIT