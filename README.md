# Harmonica demos

Eight independent terminal interfaces demonstrate small, believable uses for
Harmonica motion—not eight modes of one progress bar.

| Demo | Run |
| --- | --- |
| File transfer monitor | `go run ./cmd/file-transfer` |
| Test dashboard | `go run ./cmd/test-dashboard` |
| Shipment tracker | `go run ./cmd/shipment-tracker` |
| Agent work queue | `go run ./cmd/agent-queue` |
| Media uploader | `go run ./cmd/media-uploader` |
| Incident response | `go run ./cmd/incident-response` |
| Software rollout | `go run ./cmd/software-rollout` |
| Support escalation | `go run ./cmd/support-escalation` |

Each is interactive and simulated: it visualizes state changes locally without
copying files, running checks, buying labels, contacting agents, or uploading
data. Harmonica supplies the spring motion; Bubble Tea drives the terminal UI.

## Showcase page

`showcase/index.html` presents one VHS recording for each independent demo.
Generate the GIFs from the repository root with:

```sh
./showcase/record-all.sh
```

The helper deliberately places VHS temporary frames under `showcase/.tmp`.
That makes it work with the Snap-packaged FFmpeg on this machine, whose sandbox
cannot read the normal system `/tmp` directory.

## Silly examples

`showcase/silly/index.html` is a high-color companion page with six playful
versions of recognizable interaction patterns. It uses the same `Silly Lab`
terminal app, selected with keys `1` through `6`, and each scenario has its own
VHS tape and expanded recording.
