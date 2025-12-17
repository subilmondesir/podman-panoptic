package podman

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
	"github.com/subilmondesir/podman-panoptic/internal/core/ports"
)

// Client implémente l'interface ports.ContainerRuntime pour Podman
type Client struct {
	httpClient *http.Client
	socketPath string
	baseURL    string
}

// NewClient crée un nouveau client Podman
func NewClient(socketPath string) (*Client, error) {
	var err error
	if socketPath == "" {
		socketPath, err = detectPodmanSocket()
		if err != nil {
			return nil, fmt.Errorf("détection socket Podman: %w", err)
		}
	}

	// Vérification simple de l'existence du fichier socket
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("socket non trouvé à %s", socketPath)
	}

	client := &Client{
		httpClient: createUnixSocketClient(socketPath),
		socketPath: socketPath,
		baseURL:    "http://localhost/v1.41", // API Docker-compatible v1.41
	}

	return client, nil
}

// Ping vérifie la connectivité avec le daemon Podman
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost/_ping", nil)
	if err != nil {
		return fmt.Errorf("création requête ping: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ConnectionError{SocketPath: c.socketPath, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon répond avec status %d", resp.StatusCode)
	}

	return nil
}

// ListContainers retourne tous les conteneurs
func (c *Client) ListContainers(ctx context.Context) ([]domain.Container, error) {
	url := fmt.Sprintf("%s/containers/json?all=true", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("création requête: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ConnectionError{SocketPath: c.socketPath, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lecture réponse: %w", err)
	}

	var rawContainers []containerListResponse
	if err := json.Unmarshal(body, &rawContainers); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	containers := make([]domain.Container, 0, len(rawContainers))
	for _, raw := range rawContainers {
		containers = append(containers, raw.toDomain())
	}

	return containers, nil
}

// InspectContainer retourne les détails étendus d'un conteneur
// Cette méthode DOIT matcher la signature définie dans internal/core/ports/runtime.go
func (c *Client) InspectContainer(ctx context.Context, containerID string) (*domain.ContainerDetails, error) {
	url := fmt.Sprintf("%s/containers/%s/json", c.baseURL, containerID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("création requête: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ConnectionError{SocketPath: c.socketPath, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &NotFoundError{ContainerID: containerID}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lecture réponse: %w", err)
	}

	var raw containerInspectResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return raw.toDomain(), nil
}

// --- Helpers Internes ---

func detectPodmanSocket() (string, error) {
	// 1. Socket système
	systemSocket := "/run/podman/podman.sock"
	if fileExists(systemSocket) {
		return systemSocket, nil
	}

	// 2. Socket utilisateur (Rootless)
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		uid := os.Getuid()
		runtimeDir = fmt.Sprintf("/run/user/%d", uid)
	}

	userSocket := filepath.Join(runtimeDir, "podman", "podman.sock")
	if fileExists(userSocket) {
		return userSocket, nil
	}

	return "", fmt.Errorf("socket Podman introuvable (testé: %s, %s)", systemSocket, userSocket)
}

func createUnixSocketClient(socketPath string) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				dialer := &net.Dialer{
					Timeout: 5 * time.Second,
				}
				return dialer.DialContext(ctx, "unix", socketPath)
			},
			DisableKeepAlives: false,
			MaxIdleConns:      10,
			IdleConnTimeout:   30 * time.Second,
		},
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Vérification de l'implémentation de l'interface
var _ ports.ContainerRuntime = (*Client)(nil)
