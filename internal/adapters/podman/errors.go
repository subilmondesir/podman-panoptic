// internal/adapters/podman/errors.go
package podman

import "fmt"

// ConnectionError indique un problème de connexion au daemon
type ConnectionError struct {
	SocketPath string
	Err        error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("impossible de se connecter au daemon Podman via %s: %v", e.SocketPath, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// APIError indique une erreur retournée par l'API Podman
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Podman erreur %d: %s", e.StatusCode, e.Message)
}

// NotFoundError indique qu'un conteneur n'existe pas
type NotFoundError struct {
	ContainerID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("conteneur %s introuvable", e.ContainerID)
}
