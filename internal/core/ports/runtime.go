// internal/core/ports/runtime.go
package ports

import (
	"context"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// ContainerRuntime définit l'interface pour interagir avec un runtime de conteneurs
// (Podman, Docker, containerd, etc.)
type ContainerRuntime interface {
	// ListContainers retourne tous les conteneurs (actifs et arrêtés)
	ListContainers(ctx context.Context) ([]domain.Container, error)

	// InspectContainer retourne les détails étendus d'un conteneur
	InspectContainer(ctx context.Context, containerID string) (*domain.ContainerDetails, error)

	// Ping vérifie la connectivité avec le daemon
	Ping(ctx context.Context) error
}
