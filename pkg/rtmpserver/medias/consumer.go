package medias

type MediaConsumer interface {
	Play(batch *MediaFrameBatch)
	ID() string
	IsClosed() bool
	Close() error
}

type MediaPushConsumer interface {
	MediaConsumer
	Target() string
}
