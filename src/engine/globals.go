package engine

import "sync"

var (
	globalRenderer Renderer
	rendererMutex  sync.RWMutex
)

// SetGlobalRenderer sets the global renderer instance
func SetGlobalRenderer(renderer Renderer) {
	rendererMutex.Lock()
	defer rendererMutex.Unlock()
	globalRenderer = renderer
}

// GetGlobalRenderer returns the global renderer instance
func GetGlobalRenderer() Renderer {
	rendererMutex.RLock()
	defer rendererMutex.RUnlock()
	return globalRenderer
}
