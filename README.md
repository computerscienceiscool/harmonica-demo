# Animated progress bar

This terminal demo uses a Harmonica spring to animate progress through five
real-world scenarios. Run it from the `examples` module:

```sh
go run ./animated-progress-bar
```

Use `1` through `5` to choose a scenario, the left/right arrow keys (or
`h`/`l`) to move through its stages, `r` to reset, and `q` to quit. The
scenarios are file transfer, test suite, order shipment, agent work queue, and
media upload.

Harmonica supplies the motion between values. It does not supply a progress-bar
widget: this example draws the bar with Lip Gloss. A conventional Charm
progress-bar widget comes from the separate `bubbles/progress` package.

## Showcase page

`showcase/index.html` presents one VHS recording for each scenario. Generate
the GIFs from the `examples` module with:

```sh
vhs showcase/file-transfer.tape
vhs showcase/test-suite.tape
vhs showcase/order-shipment.tape
vhs showcase/agent-work-queue.tape
vhs showcase/media-upload.tape
```
