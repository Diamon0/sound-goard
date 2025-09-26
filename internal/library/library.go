package library

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Diamon0/sound-goard/internal/library/playback"
)

// TODO: Move ENV variables into a config file

// The maximum duration, in milliseconds, to playback media for
var PLAYBACK_TIMEOUT int = 120 * 1000

func init() {
	if playbackTimeout := os.Getenv("PLAYBACK_TIMEOUT"); playbackTimeout != "" {
		if newTimeout, err := strconv.Atoi(playbackTimeout); err != nil {
			fmt.Println("Failed to parse PLAYBACK_TIMEOUT fron env")
		} else {
			PLAYBACK_TIMEOUT = newTimeout
		}
	}
}

// Defines a collection of playable media, as well as its state
type Library struct {
	Media []playback.Playable
	CurrentlyPlaying []*playback.Playing
}

// Defines a media file
type MediaFile struct {
	// The media name
	Name string

	// The path to the media file
	FilePath string

	// The length, in milliseconds, of the media
	Length int
}

func (mf *MediaFile) Start(ctx context.Context) (playback.Playing, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(PLAYBACK_TIMEOUT))

	cmd := exec.CommandContext(ctx, "ffplay", "-nodisp", "-v", "quiet", "-autoexit", mf.FilePath)

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, errors.Join(fmt.Errorf("Failed to play media file %s: ", mf.FilePath), err)
	}

	return &PlaybackInstance{
		Name: mf.Name,
		Ctx: &ctx,
		Cancel: cancel,
		Cmd: cmd,
	}, nil
}

// Defines a Playable instance for playback
type PlaybackInstance struct {
	// The name of the playing media
	Name string

	// The context of the playback instance
	Ctx *context.Context

	// The cancel function for the playback instance's context
	Cancel context.CancelFunc

	// The command being used for playback
	Cmd *exec.Cmd
}

// INFO: Currently, the error return goes unused
func (pi *PlaybackInstance) Stop() error {
	pi.Cancel()
	return nil
}
