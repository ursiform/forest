// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import "io"

type Populater interface {
	Populate(body io.ReadCloser) error
}
