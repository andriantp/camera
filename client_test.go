package camera

import (
	"testing"
	"time"
)

func TestNewDefaults(t *testing.T) {
	c := New(Config{})

	if got := c.HTTP().Timeout; got != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", got)
	}
	if got := c.HTTPDownload().Timeout; got != 5*time.Minute {
		t.Errorf("DownloadTimeout = %v, want 5m", got)
	}
	if got := c.MoveTimeout(); got != 5*time.Second {
		t.Errorf("MoveTimeout = %v, want 5s", got)
	}
}

func TestNewCustomConfig(t *testing.T) {
	c := New(Config{
		Timeout:         10 * time.Second,
		DownloadTimeout: time.Minute,
		MoveTimeout:     20 * time.Second,
	})

	if got := c.HTTP().Timeout; got != 10*time.Second {
		t.Errorf("Timeout = %v", got)
	}
	if got := c.HTTPDownload().Timeout; got != time.Minute {
		t.Errorf("DownloadTimeout = %v", got)
	}
	if got := c.MoveTimeout(); got != 20*time.Second {
		t.Errorf("MoveTimeout = %v", got)
	}
	if c.HTTP().Transport != c.HTTPDownload().Transport {
		t.Error("regular and download clients should share one transport")
	}
}
