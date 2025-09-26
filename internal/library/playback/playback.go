package playback

// Defines a Playable entity
type Playable interface {
	Start() (*Playing, error)
}

// Defines an entity currently being played
type Playing interface {
	Stop() error
}
