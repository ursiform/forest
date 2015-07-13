package forest

import "io"

type Populater interface {
	Populate(body io.ReadCloser) error
}
