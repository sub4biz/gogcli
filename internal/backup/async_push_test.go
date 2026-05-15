package backup

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitAsyncPushRetryCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitAsyncPushRetry(ctx, time.Hour)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestAsyncPusherProgressCanBeUpdated(t *testing.T) {
	var got []string
	p := &asyncRepoPusher{
		progress: func(format string, _ ...any) {
			got = append(got, "old:"+format)
		},
	}

	p.setProgress(func(format string, _ ...any) {
		got = append(got, "new:"+format)
	})
	p.progressf("event")

	if len(got) != 1 || got[0] != "new:event" {
		t.Fatalf("unexpected progress calls: %#v", got)
	}
}
