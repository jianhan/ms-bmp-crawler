package outputs

import (
	"context"

	"github.com/jianhan/ms-bmp-crawler/crawlers"
)

type OutputWriter interface {
	Output(ctx context.Context, crawler crawlers.Crawler) error
}
