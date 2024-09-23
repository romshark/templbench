package templates

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

func BenchmarkPageHelloWorld(b *testing.B) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	buf.Grow(4096)
	ctx = templ.InitializeContext(ctx)
	b.ResetTimer()

	for range b.N {
		if err := RenderHelloWorld(ctx, buf, "title", "msg"); err != nil {
			panic(err)
		}
		buf.Reset()
	}
}
