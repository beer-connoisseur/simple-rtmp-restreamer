module github.com/kbats183/simple-rtmp-restreamer

go 1.27.0

toolchain go1.27.3

require (
	github.com/go-chi/chi/v5 v5.5.0
	github.com/yapingcat/gomedia v0.0.0-20240906162731-17feea57090c
)

replace github.com/yapingcat/gomedia => github.com/kbats183/gomedia v0.0.0-20250817114334-50ae796beb83
