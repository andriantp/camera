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

| Area | Methods |
|---|---|
| Device | GetDeviceInfo, Reboot |
| Capabilities | GetCapabilities, GetSystemCapabilities, GetStreamingCapabilities, GetPTZCapabilities, GetEventCapabilities, GetImageCapabilities, GetAnalyticsCapabilities |
| Video | GetInputs, GetStreaming, GetSnapshot |
| PTZ | GetStatus, GetAbsolute, SetAbsolute, SetZoom, RelativeMove, GotoHome, GetPresets, GotoPreset |
| Recording | GetContentCapabilities, GetStorage, ContentSearch, ContentDownload |
| Events | GetHTTPHosts, SetHTTPHost, GetMotionDetection, SetMotionSensitivity |
| Thermal | SetThermometry, GetThermometryRulesTemperatureInfo |

## Installation

```bash
go get github.com/andriantp/camera
```

Requires Go 1.23 or later.

## Quick Start

### ONVIF

```go
client := camera.New(camera.Config{})

onvifClient := onvif.New(client)

capabilities, err := onvifClient.GetCapabilities(
    context.Background(),
    onvif.Camera{
        Credential: camera.Credential{
            Host:     "192.168.1.64",
            Username: "admin",
            Password: "...",
        },
        Profile: "profile_1",
    },
)
if err != nil {
    log.Fatal(err)
}
```

### Hikvision ISAPI

```go
client := camera.New(camera.Config{})

isapiClient := isapi.New(client)

capabilities, err := isapiClient.GetCapabilities(
    context.Background(),
    isapi.Camera{
        Credential: camera.Credential{
            Host:     "192.168.1.64",
            Username: "admin",
            Password: "...",
        },
        Channel: 1,
    },
)
if err != nil {
    log.Fatal(err)
}
```

Create one `camera.Client` and one driver client per application and reuse
them: they are safe for concurrent use and cache per-camera state (see
[Behavior](#behavior)).

## Configuration

```go
client := camera.New(camera.Config{
    Timeout:         30 * time.Second, // HTTP timeout per request (default 30s)
    DownloadTimeout: 5 * time.Minute,  // HTTP timeout for recording downloads (default 5m)
    MoveTimeout:     10 * time.Second, // how long SetAbsolute/SetZoom wait for the camera to reach the target (default 5s)
    TLSConfig:       nil,              // TLS settings for https:// hosts (default: verify against system roots)
})
```

All fields are optional. A shorter `context` deadline always takes precedence.
When `MoveTimeout` elapses, the error wraps `isapi.ErrMoveTimeout`:

```go
if errors.Is(err, isapi.ErrMoveTimeout) {
    // camera did not reach the target in time
}
```

### Camera host

`Credential.Host` accepts:

| Host | Connects to |
|---|---|
| `192.168.1.64` | `http://192.168.1.64` |
| `192.168.1.64:8080` | `http://192.168.1.64:8080` |
| `https://192.168.1.64` | `https://192.168.1.64` (ISAPI) |
| `https://192.168.1.64:8443` | `https://192.168.1.64:8443` (ISAPI) |

Schemes other than `http` and `https` are rejected.

## HTTPS

Prefix the camera host with `https://` to use HTTPS for the ISAPI driver:

```go
isapi.Camera{
    Credential: camera.Credential{
        Host:     "https://192.168.1.64", // or "https://192.168.1.64:8443"
        Username: "admin",
        Password: "...",
    },
    Channel: 1,
}
```

Cameras usually ship a self-signed certificate, which is rejected by default.
Trust it explicitly with `Config.TLSConfig`:

```go
pem, err := os.ReadFile("camera.crt") // exported from the camera web UI
if err != nil {
    log.Fatal(err)
}
roots := x509.NewCertPool()
roots.AppendCertsFromPEM(pem)

client := camera.New(camera.Config{
    TLSConfig: &tls.Config{RootCAs: roots},
})
```

If verification fails with a hostname error, the certificate does not list
the address you connect to: regenerate it on the camera with the right IP or
hostname, or set `ServerName` in `TLSConfig` to a name the certificate covers.

`&tls.Config{InsecureSkipVerify: true}` also works, but disables certificate
checks and exposes credentials to man-in-the-middle attacks; use it only on a
trusted network.

The ONVIF driver always uses HTTP (the underlying library hardcodes it). An
`https://` host is accepted so one `Credential` can be shared with ISAPI: the
scheme and HTTPS port are dropped and ONVIF connects to the default HTTP port.
If ONVIF listens on a non-default port, give the ONVIF client its own
`Credential` with `"host:port"`.

## Behavior

- **Authentication** — ISAPI uses HTTP digest auth. The digest challenge is
  cached per credential, so after the first request each call needs a single
  round trip.
- **ONVIF device cache** — service endpoints are discovered once per camera
  and reused. Calls to the same camera run one at a time; different cameras
  run in parallel. Create a new `onvif.Client` if a camera's endpoints change
  (e.g. after a firmware update).
- **Context** — every method honors `ctx` cancellation and deadlines,
  including ONVIF calls.
- **`isapi.GetCapabilities`** — fetches all capability groups concurrently.
  Groups the camera does not support are left `nil`; an error is returned only
  if every group fails or `ctx` is done.
- **Waiting moves** — `SetAbsolute` and `SetZoom` poll until the camera
  reaches the target (within `MoveTimeout`). `RelativeMove` waits for the
  given duration before returning the new status.
- **Streams** — `GetSnapshot` and `ContentDownload` return an `io.ReadCloser`.
  Read it to the end and close it so the connection can be reused. The
  returned size is `-1` when the camera does not send `Content-Length`.

## Examples

Runnable examples live in [`examples/`](examples). They read the camera
connection from `examples/common/.env`:

```bash
cp examples/common/.env.example examples/common/.env
# edit examples/common/.env
```

| Variable | Required | Description |
|---|---|---|
| `Host` | yes | Camera address, see [Camera host](#camera-host) |
| `User` | yes | Camera username |
| `Password` | yes | Camera password (quote it if it contains spaces or ` #`) |
| `CACert` | no | Path to the camera certificate (PEM) to trust for `https://` |
| `Insecure` | no | `true` to skip TLS certificate verification |
| `Channel` | no | ISAPI channel (default `1`) |
| `Profile` | no | ONVIF profile token (default `profile_1`) |

`.env` is git-ignored; only `.env.example` is committed.

Run an example from the repository root:

```bash
go run ./examples/onvif/status
go run ./examples/isapi/channels INPUTS
```

| Example | Arguments |
|---|---|
| `onvif/capability` | — |
| `onvif/configuration` | — |
| `onvif/profiles` | — |
| `onvif/status` | — |
| `onvif/relative_move` | `LEFT` \| `RIGHT` \| `UP` \| `DOWN` |
| `isapi/capability` | `DEVICE` \| `ALL` \| `SYSTEM` \| `STREAMING` \| `PTZ` \| `EVENT` \| `IMAGE` \| `ANALYTIC` |
| `isapi/channels` | `INPUTS` \| `STREAMING` |
| `isapi/status` | — |
| `isapi/home` | — |
| `isapi/preset` | `GET` \| `GOTO <id>` |
| `isapi/relative_move` | `LEFT` \| `RIGHT` \| `UP` \| `DOWN` |
| `isapi/zoom` | `RANGE` \| `CURRENT` \| `GETABSOLUTE` \| `ZOOM <value>` \| `SETABSOLUTE` |
| `isapi/content` | `CAPABILITIES` \| `STORAGE` \| `SEARCH` \| `DOWNLOAD` \| `SNAPSHOT` |

`isapi/zoom SETABSOLUTE` and `isapi/content DOWNLOAD` use values set in their
`main.go`; edit them for your camera. Snapshots and downloads are written to
the current directory.

## Testing

```bash
go test ./... -race
```

The tests run against in-process fake cameras (`httptest`), so no device is
needed.

## Project Structure

```text
.
├── client.go          # camera.Client and Config
├── driver/
│   ├── isapi/         # Hikvision ISAPI driver
│   │   ├── xml/       # ISAPI XML parsing
│   │   └── json/      # ISAPI JSON parsing
│   └── onvif/         # ONVIF driver
│       └── xml/       # ONVIF SOAP parsing
├── model/             # Shared public models
├── examples/          # Runnable examples
└── README.md
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
| ONVIF | ✅ (HTTP) |
| Hikvision ISAPI | ✅ (HTTP / HTTPS) |

## License

MIT
