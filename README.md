# 👁️ PANOPTIC

🇬🇧**The all-seeing eye for Podman security** • **L'œil omniscient pour la sécurité Podman**  
🇫🇷**Next-gen container audit system • Système d'audit nouvelle génération**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Podman](https://img.shields.io/badge/Podman-Native-892CA0?style=for-the-badge&logo=podman)](https://podman.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/badge/Release-v1.0.0-gold?style=for-the-badge)](https://github.com/subilmondesir/podman-panoptic/releases)

---

> **Philosophy:** *"See Everything. Secure Everywhere."*  
> **Philosophie :** *"Voir tout. Sécuriser partout."*

🇬🇧**PANOPTIC** is a next-generation security audit tool for **Podman** environments, written in **Go**. It performs deep inspection via the Podman Socket API, integrates **Trivy** for CVE scanning, and delivers results through an immersive Terminal UI or comprehensive HTML reports.

🇫🇷**PANOPTIC** est un outil d'audit de sécurité nouvelle génération pour les environnements **Podman**, écrit en **Go**. Il effectue une inspection approfondie via l'API Socket Podman, intègre **Trivy** pour le scanning CVE, et génère des rapports via une interface terminal immersive ou des rapports HTML complets.

---

## ⚡ Features • Fonctionnalités 🇬🇧🇫🇷

| Feature | Description (EN/FR) |
|---------|---------------------|
| **🔭 Deep Inspection** | Direct Podman Socket API communication (Rootless & Rootful) <br> Communication directe avec l'API Socket Podman (Sans root & Avec root) |
| **🛡️ Trivy Core** | **Active:** Embedded CVE scanning engine <br> **Actif :** Moteur de scanning CVE intégré |
| **⚡ Fast Scanning** | Multi-threaded via Goroutines (< 10s scan time) <br> Multi-threadé via Goroutines (scan en < 10s) |
| **🎮 Interactive TUI** | **Active:** Real-time dashboard powered by Bubble Tea <br> **Actif :** Tableau de bord temps réel propulsé par Bubble Tea |
| **📄 Smart Reports** | JSON (CI/CD) or Modern HTML5 with Risk Scoring <br> JSON (CI/CD) ou HTML5 moderne avec Score de Risque |
| **🏗️ Clean Code** | Hexagonal Architecture for maintainability <br> Architecture Hexagonale pour la maintenabilité |

---

## 🚀 Quick Start • Démarrage Rapide 🇬🇧🇫🇷

### Prerequisites • Prérequis

- **Linux/macOS** (Windows via WSL2)
- **Podman** installed and running • installé et actif
- **Trivy** installed (optional for CVE features) • installé (optionnel pour les CVE)
- **Go 1.22+** (if building from source) • (si compilation depuis source)

---

### Installation

#### Option 1: Download Binary • Télécharger le binaire (Recommandé)
```bash
# Download latest release • Télécharger la dernière version
wget https://github.com/subilmondesir/podman-panoptic/releases/download/v1.0.0/panoptic-linux-amd64

# Make executable • Rendre exécutable
chmod +x panoptic-linux-amd64

# Install to PATH • Installer dans le PATH
sudo mv panoptic-linux-amd64 /usr/local/bin/panoptic
```

#### Option 2: Build from Source • Compiler depuis la source
```bash
# Clone repository • Cloner le dépôt
git clone https://github.com/subilmondesir/podman-panoptic.git
cd panoptic

# Download dependencies • Télécharger les dépendances
go mod tidy

# Build Gold Master binary • Compiler le binaire Gold Master
make build

# Install (optional) • Installer (optionnel)
sudo make install
```

---

## 🎮 Usage • Utilisation 🇬🇧🇫🇷

### 1. Start Podman Socket • Démarrer le Socket Podman

PANOPTIC requires the Podman socket to be active.  
PANOPTIC nécessite que le socket Podman soit actif.
```bash
# Rootless mode • Mode sans root
podman system service --time=0 &

# Rootful mode (with sudo) • Mode root (avec sudo)
sudo podman system service --time=0 &
```

---

### 2. Run Security Audit • Lancer l'Audit de Sécurité

**Interactive Mode (TUI) • Mode Interactif (TUI):**
```bash
panoptic scan
# Real-time progress dashboard • Tableau de bord en temps réel
```

**HTML Report (Rich UI) • Rapport HTML (Interface Riche):**
```bash
panoptic scan --format html --output security-report.html
# Professional report for teams • Rapport professionnel pour équipes
```

**JSON for CI/CD • JSON pour CI/CD:**
```bash
panoptic scan --format json --output audit.json
# Machine-readable output • Sortie lisible par machine
```

**Custom Timeout • Timeout Personnalisé:**
```bash
panoptic scan --timeout 300
# For slow networks or large environments • Pour réseaux lents ou grands environnements
```

---

## 📊 Example Output • Exemple de Sortie 🇬🇧🇫🇷

### Terminal Report • Rapport Terminal
```
======================================================================
👁️  PANOPTIC - SECURITY AUDIT REPORT • RAPPORT D'AUDIT DE SÉCURITÉ
======================================================================

Date:         2025-12-15 03:14:02
Version:      1.0.0
Duration:     8.6s • Durée: 8.6s
Containers:   3 • Conteneurs: 3

----------------------------------------------------------------------
📊 EXECUTIVE SUMMARY • RÉSUMÉ EXÉCUTIF
----------------------------------------------------------------------
CVE Vulnerabilities:  0 (Critical: 0, High: 0)
Misconfigurations:    3
Privileged Containers: 1
Risk Score:           3.0/100 • Score de Risque: 3.0/100

----------------------------------------------------------------------
🛡️  SECURITY FINDINGS • DÉCOUVERTES DE SÉCURITÉ
----------------------------------------------------------------------

[1] 🔴 CRITICAL • CRITIQUE
    PANOPTIC-003 - Secrets in Environment Variables
                   Secrets dans les Variables d'Environnement
    
    Resource:     secret-leak
    Detected:     AWS_SECRET_KEY
    
    💡 Remediation:
       EN: Use secrets managers (Podman secrets, Vault)
       FR: Utiliser des gestionnaires de secrets (Podman secrets, Vault)
```

---

## 🧪 Capabilities Matrix • Matrice des Capacités (v1.0.0) 🇬🇧🇫🇷
### ✅ Core Security Engine • Moteur de Sécurité Principal (Active • Actif)

| Capability | Description |
|------------|-------------|
| **CVE Analysis** | Deep image scanning via embedded Trivy integration <br> Scan approfondi d'images via intégration Trivy embarquée |
| **Privileged Detection** | Identifies containers with root capabilities (`--privileged`) <br> Identifie les conteneurs avec capacités root (`--privileged`) |
| **Sensitive Mounts** | Detects risky volume bindings (`/etc`, `/proc`, `/sys`) <br> Détecte les montages risqués (`/etc`, `/proc`, `/sys`) |
| **Secret Leakage** | Heuristic analysis of env vars (`AWS_KEY`, `PASSWORD`, `TOKEN`) <br> Analyse heuristique des variables d'env (`AWS_KEY`, `PASSWORD`, `TOKEN`) |
| **Network Exposure** | Flags Host Network mode usage (`--net=host`) <br> Signale l'usage du mode réseau hôte (`--net=host`) |

### 🚀 User Experience • Expérience Utilisateur (Active • Actif)

- **Interactive TUI:** Real-time progress tracking with Bubble Tea  
  **TUI Interactif :** Suivi de progression en temps réel avec Bubble Tea

- **Smart Reporting:** HTML5 reports with Risk Score calculation  
  **Rapports Intelligents :** Rapports HTML5 avec calcul du Score de Risque

### 🔮 Roadmap (v1.1 • Feuille de Route)

- [ ] **Auto-remediation:** One-click fixes for common issues  
      **Auto-remédiation :** Corrections en un clic pour problèmes courants

- [ ] **Live Watch Mode:** Daemon mode for container spawn monitoring  
      **Mode Surveillance Live :** Mode daemon pour monitoring des spawns

- [ ] **Remote Scanning:** Audit remote Podman instances via SSH  
      **Scan Distant :** Audit d'instances Podman distantes via SSH

---

## 🏗️ Architecture 🇬🇧🇫🇷

PANOPTIC follows strict **Hexagonal Architecture** (Ports & Adapters) for modularity.  
PANOPTIC suit une **Architecture Hexagonale** stricte (Ports & Adapters) pour la modularité.
```
internal/
├── core/              # Business logic • Logique métier
│   ├── domain/        # Entities (Container, Vulnerability)
│   ├── ports/         # Interfaces (Runtime, Scanner, Reporter)
│   └── services/      # Orchestration (AuditService)
├── adapters/          # Infrastructure implementations
│   ├── podman/        # Podman HTTP API client
│   ├── trivy/         # Trivy CLI Wrapper
│   └── system/        # System compliance checks
└── ui/                # User interfaces
    ├── cli/           # Cobra command system
    ├── tui/           # Bubble Tea interactive UI
    └── output/        # Report generators (HTML/JSON)
```

---

## 🔧 Configuration 🇬🇧🇫🇷

### Config File • Fichier de Configuration (Optional • Optionnel)

Create `~/.panoptic.yaml`:
```yaml
output:
  format: html        # text, json, html
  
scan:
  timeout: 30         # seconds • secondes
  
verbose: false        # detailed logs • logs détaillés
```

### Environment Variables • Variables d'Environnement
```bash
export PANOPTIC_OUTPUT_FORMAT=json
export PANOPTIC_VERBOSE=true
```

---

## 🧪 Performance Metrics • Métriques de Performance 🇬🇧🇫🇷

*Stress-tested on Kali Linux, Debian, Fedora, AlmaLinux, Ubuntu with Podman 5.4.2*  
*Testé en conditions de stress sur Kali Linux, Debian, Fedora, AlmaLinux, Ubuntu avec Podman 5.4.2*

| Metric • Métrique | Result • Résultat | Verdict |
|-------------------|-------------------|---------|
| **Binary Size** | ~6MB (static) | ✅ Optimized • Optimisé |
| **Scan Speed** | **8.6s** (3 containers • 3 conteneurs) | ⚡ Blazing Fast • Ultra-rapide |
| **CVE Detection** | Active (Trivy Core) | 🎯 Operational • Opérationnel |
| **Secret Detection** | AWS Keys detected • Clés AWS détectées | 🔴 Critical Alert • Alerte Critique |
| **Memory Usage** | <50MB during scan • <50MB pendant scan | ✅ Efficient • Efficace |

---

## 🤝 Contributing • Contribuer 🇬🇧🇫🇷

Contributions are welcome! Please maintain the hexagonal architecture pattern.  
Les contributions sont bienvenues ! Merci de maintenir le pattern d'architecture hexagonale.

**Guidelines:**
- **Core logic** → `internal/core`
- **Infrastructure** → `internal/adapters`
- **Tests required** • Tests requis for new features

---

## 📜 License • Licence 🇬🇧🇫🇷

MIT License - see [LICENSE](LICENSE) file for details.  
Licence MIT - voir le fichier [LICENSE](LICENSE) pour détails.

---

## 🙏 Acknowledgments • Remerciements 

- **Podman Project** for the amazing container runtime  
- **Aqua Security** for Trivy CVE scanner  
- **Charm.sh** for Bubble Tea TUI framework  
- **Go Community** for excellent tooling

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/subilmondesir/podman-panoptic/issues)
- **Documentation:** [Wiki](https://github.com/subilmondesir/podman-panoptic/wiki)
- **Author:** **XarKEzion** [@subilmondesir](https://github.com/subilmondesir)

---

<div align="center">

**Built with precision and passion**  
**Construit avec précision et passion**

*L'Artisan du Code Horloger* 🕰️

**© 2025 PANOPTIC Project**

</div>
