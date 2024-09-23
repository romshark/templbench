package templates

import (
	"context"
	"io"
)

//go:generate go run github.com/a-h/templ/cmd/templ generate

func RenderHelloWorld(ctx context.Context, w io.Writer, title, msg string) error {
	return pageHelloWorld(title, msg).Render(ctx, w)
}
