# Harmonica demos

Five independent terminal interfaces demonstrate small, believable uses for
Harmonica motion—not five modes of one progress bar.

| Demo | Run |
| --- | --- |
| File transfer monitor | `go run ./cmd/file-transfer` |
| Test dashboard | `go run ./cmd/test-dashboard` |
| Shipment tracker | `go run ./cmd/shipment-tracker` |
| Agent work queue | `go run ./cmd/agent-queue` |
| Media uploader | `go run ./cmd/media-uploader` |

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
