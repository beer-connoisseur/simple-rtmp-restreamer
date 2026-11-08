# Monitoring

## Metrics

Prometheus metrics are served at `GET /metrics` on the HTTP port (`6070`).
If `BASIC_AUTH_USER` / `BASIC_AUTH_PASS` are set, the endpoint requires the same basic auth as the API.

Besides the standard Go runtime (`go_*`) and process (`process_*`) metrics, the restreamer exports:

| Metric | Type | Labels | Description |
|---|---|---|---|
| `restreamer_rtmp_sessions` | gauge | | Open incoming RTMP connections (publishers and players) |
| `restreamer_ingest_active` | gauge | `stream` | 1 while a publisher is sending the stream |
| `restreamer_ingest_frames_total` | counter | `stream`, `kind`, `codec` | Frames received from the publisher |
| `restreamer_ingest_bytes_total` | counter | `stream`, `kind`, `codec` | Payload bytes received from the publisher |
| `restreamer_ingest_keyframes_total` | counter | `stream` | Video keyframes received |
| `restreamer_stream_dropped_batches_total` | counter | `stream` | Batches dropped: stream dispatch buffer full |
| `restreamer_push_connected` | gauge | `stream`, `target` | 1 while the push target accepts the publish |
| `restreamer_push_connect_failures_total` | counter | `stream`, `target` | Push sessions that ended with an error |
| `restreamer_push_errors_total` | counter | `stream`, `target`, `code` | RTMP error statuses reported by the target |
| `restreamer_push_bytes_total` | counter | `stream`, `target` | Payload bytes sent to the target |
| `restreamer_push_queue_batches` | gauge | `stream`, `target` | Batches waiting to be sent to the target |
| `restreamer_push_dropped_batches_total` | counter | `stream`, `target` | Batches dropped: queue overflow while the target is connected |
| `restreamer_pull_consumers` | gauge | `stream` | Players watching the stream from the restreamer |
| `restreamer_pull_dropped_batches_total` | counter | `stream` | Batches dropped: player queue overflow |

The `target` label is the target name from the API (or the target host if the name is empty),
never the full URL, so stream keys do not leak into metrics.