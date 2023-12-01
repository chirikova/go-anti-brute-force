package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	defer goleak.VerifyNone(t)

	type fields struct {
		limit    int64
		interval time.Duration
	}
	type args struct {
		count int
		key   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		allow  bool
	}{
		{
			name: "Request allowed",
			fields: fields{
				interval: time.Minute,
				limit:    5,
			},
			args: args{
				key:   "testkey",
				count: 2,
			},
			allow: true,
		},
		{
			name: "Request denied",
			fields: fields{
				interval: time.Minute,
				limit:    5,
			},
			args: args{
				key:   "testkey",
				count: 5,
			},
			allow: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			limiter := NewSlidingWindowLimiter(ctx, tt.fields.interval, tt.fields.limit)

			for i := 0; i < tt.args.count; i++ {
				_ = limiter.Allow(tt.args.key)
			}

			ok := limiter.Allow(tt.args.key)
			cancel()

			require.Equal(t, ok, tt.allow)
		})
	}
}

func TestSlidingWindowLimiter_Clean(t *testing.T) {
	defer goleak.VerifyNone(t)

	type fields struct {
		limit    int64
		interval time.Duration
	}
	type args struct {
		count int
		keys  []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Clean is successful",
			fields: fields{
				interval: time.Second * 1,
				limit:    5,
			},
			args: args{
				keys:  []string{"testkey1", "testkey2"},
				count: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			limiter := NewSlidingWindowLimiter(ctx, tt.fields.interval, tt.fields.limit)

			for _, key := range tt.args.keys {
				for i := 0; i < tt.args.count; i++ {
					_ = limiter.Allow(key)
				}

				ok := limiter.Allow(key)
				require.False(t, ok)
			}
			time.Sleep(tt.fields.interval)

			limiter.Clean()

			for _, key := range tt.args.keys {
				ok := limiter.Allow(key)
				require.True(t, ok)
			}
			cancel()
		})
	}
}

func TestSlidingWindowLimiterReset(t *testing.T) {
	defer goleak.VerifyNone(t)

	type fields struct {
		limit    int64
		interval time.Duration
	}
	type args struct {
		count int
		key   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Reset is successful",
			fields: fields{
				interval: time.Minute,
				limit:    5,
			},
			args: args{
				key:   "testkey",
				count: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			limiter := NewSlidingWindowLimiter(ctx, tt.fields.interval, tt.fields.limit)

			for i := 0; i < tt.args.count; i++ {
				_ = limiter.Allow(tt.args.key)
			}

			ok := limiter.Allow(tt.args.key)
			require.False(t, ok)

			limiter.Reset(tt.args.key)

			ok = limiter.Allow(tt.args.key)
			require.True(t, ok)

			cancel()
		})
	}
}
