// internal/adapters/podman/types.go
package podman

import (
	"time"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// containerListResponse structure de la réponse API /containers/json
type containerListResponse struct {
	ID      string            `json:"Id"`
	Names   []string          `json:"Names"`
	Image   string            `json:"Image"`
	State   string            `json:"State"`
	Status  string            `json:"Status"`
	Created int64             `json:"Created"`
	Labels  map[string]string `json:"Labels"`
}

// toDomain convertit la réponse API en domain.Container
func (r containerListResponse) toDomain() domain.Container {
	name := "<unnamed>"
	if len(r.Names) > 0 {
		name = r.Names[0]
		// Supprimer le "/" préfixe Podman
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}
	}

	// Truncate ID à 12 chars (style Docker)
	id := r.ID
	if len(id) > 12 {
		id = id[:12]
	}

	return domain.Container{
		ID:      id,
		Name:    name,
		Image:   r.Image,
		State:   domain.ContainerState(r.State),
		Status:  r.Status,
		Created: time.Unix(r.Created, 0),
		Labels:  r.Labels,
	}
}

// containerInspectResponse structure de la réponse API /containers/{id}/json
type containerInspectResponse struct {
	ID         string                  `json:"Id"`
	Name       string                  `json:"Name"`
	Image      string                  `json:"Image"`
	State      containerStateResponse  `json:"State"`
	Created    string                  `json:"Created"`
	Config     containerConfigResponse `json:"Config"`
	HostConfig hostConfigResponse      `json:"HostConfig"`
	Mounts     []mountResponse         `json:"Mounts"`
}

type containerStateResponse struct {
	Status     string `json:"Status"`
	Running    bool   `json:"Running"`
	Paused     bool   `json:"Paused"`
	Restarting bool   `json:"Restarting"`
	Pid        int    `json:"Pid"`
}

type containerConfigResponse struct {
	Labels map[string]string `json:"Labels"`
	Env    []string          `json:"Env"`
}

type hostConfigResponse struct {
	Privileged  bool   `json:"Privileged"`
	NetworkMode string `json:"NetworkMode"`
}

type mountResponse struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
}

// toDomain convertit la réponse inspect en domain.ContainerDetails
func (r containerInspectResponse) toDomain() *domain.ContainerDetails {
	// Parse de la date de création
	created, _ := time.Parse(time.RFC3339, r.Created)

	// Construction de l'état
	var state domain.ContainerState
	switch {
	case r.State.Running:
		state = domain.StateRunning
	case r.State.Paused:
		state = domain.StatePaused
	case r.State.Restarting:
		state = domain.StateRestarting
	default:
		state = domain.ContainerState(r.State.Status)
	}

	// Nom sans le "/" préfixe
	name := r.Name
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	// ID truncate
	id := r.ID
	if len(id) > 12 {
		id = id[:12]
	}

	// Conversion des mounts
	mounts := make([]domain.Mount, len(r.Mounts))
	for i, m := range r.Mounts {
		mode := "ro"
		if m.RW {
			mode = "rw"
		}
		mounts[i] = domain.Mount{
			Type:        m.Type,
			Source:      m.Source,
			Destination: m.Destination,
			Mode:        mode,
		}
	}

	// Parsing des variables d'environnement
	envVars := make(map[string]string)
	for _, env := range r.Config.Env {
		// Format: "KEY=VALUE"
		for i, char := range env {
			if char == '=' {
				key := env[:i]
				value := env[i+1:]
				envVars[key] = value
				break
			}
		}
	}

	return &domain.ContainerDetails{
		Container: domain.Container{
			ID:      id,
			Name:    name,
			Image:   r.Image,
			State:   state,
			Status:  r.State.Status,
			Created: created,
			Labels:  r.Config.Labels,
		},
		Privileged:      r.HostConfig.Privileged,
		Mounts:          mounts,
		NetworkMode:     r.HostConfig.NetworkMode,
		EnvironmentVars: envVars,
		PID:             r.State.Pid,
	}
}
