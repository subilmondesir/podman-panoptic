// internal/core/domain/container.go
package domain

import "time"

// Container représente un conteneur en cours d'exécution ou arrêté
type Container struct {
	ID      string            // Identifiant unique (12 premiers chars du hash)
	Name    string            // Nom humainement lisible
	Image   string            // Image source (ex: nginx:alpine)
	State   ContainerState    // État actuel (running, exited, etc.)
	Status  string            // Description textuelle de l'état
	Created time.Time         // Date de création
	Labels  map[string]string // Métadonnées clé-valeur
}

// ContainerState représente l'état d'un conteneur
type ContainerState string

const (
	StateRunning    ContainerState = "running"
	StatePaused     ContainerState = "paused"
	StateRestarting ContainerState = "restarting"
	StateExited     ContainerState = "exited"
	StateCreated    ContainerState = "created"
	StateDead       ContainerState = "dead"
)

// ContainerDetails contient les métadonnées étendues d'un conteneur
type ContainerDetails struct {
	Container            // Embedding des infos de base
	Privileged      bool // Accès root complet au host
	Mounts          []Mount
	NetworkMode     string
	EnvironmentVars map[string]string
	PID             int
}

// Mount représente un point de montage filesystem
type Mount struct {
	Type        string // bind, volume, tmpfs
	Source      string // Chemin source
	Destination string // Chemin dans le conteneur
	Mode        string // ro, rw
}

// IsRunning vérifie si le conteneur est actif
func (c Container) IsRunning() bool {
	return c.State == StateRunning
}

// IsPrivileged vérifie si le conteneur a des privilèges élevés
func (cd ContainerDetails) IsPrivileged() bool {
	return cd.Privileged
}

// HasSensitiveMounts détecte les montages sensibles (/etc, /var, /root)
func (cd ContainerDetails) HasSensitiveMounts() bool {
	sensitivePaths := []string{"/etc", "/var", "/root", "/sys", "/proc"}

	for _, mount := range cd.Mounts {
		for _, sensitive := range sensitivePaths {
			if mount.Source == sensitive || mount.Destination == sensitive {
				return true
			}
		}
	}
	return false
}
