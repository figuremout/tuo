package all

import (
	// Blank imports for plugins to register themselves
	// init func will be run automatically when package is imported
	_ "github.com/githubzjm/tuo/internal/pkg/metrics/cpu"
	_ "github.com/githubzjm/tuo/internal/pkg/metrics/heartbeat"
	_ "github.com/githubzjm/tuo/internal/pkg/metrics/host"
	_ "github.com/githubzjm/tuo/internal/pkg/metrics/mem"
)
