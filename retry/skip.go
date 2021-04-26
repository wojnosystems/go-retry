package retry

import (
	"github.com/wojnosystems/go-retry/core"
)

// Skip is a test helping placeholder that will not call the method even once
// You can use this to disable some block of logic guarded with a Retry
type Skip struct {
}

func (s *Skip) Retry(_ core.CallbackFunc) (err error) {
	return nil
}
