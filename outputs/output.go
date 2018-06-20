package outputs

import "github.com/jianhan/ms-bmp-crawler/crawlers"

type OutputWriter interface {
	Output(crawler crawlers.Crawler) error
}
