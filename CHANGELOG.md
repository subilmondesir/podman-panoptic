# Changelog

All notable changes to PANOPTIC are documented here.  
Tous les changements notables de PANOPTIC sont documentés ici.

Format based on [Keep a Changelog](https://keepachangelog.com/).  
Format basé sur [Keep a Changelog](https://keepachangelog.com/).

---

## [1.0.0] - 2025-12-15

### 🎉 Gold Master Release • Version Gold Master

First production-ready release of PANOPTIC.  
Première version production-ready de PANOPTIC.

#### ✨ Added • Ajouté

**Core Architecture • Architecture Principale**
- Hexagonal design with clean separation of concerns  
  Architecture hexagonale avec séparation propre des responsabilités
- Domain-driven design for business logic  
  Conception orientée domaine pour la logique métier

**Podman Integration • Intégration Podman**
- Native HTTP API communication via Unix socket  
  Communication API HTTP native via socket Unix
- Support for both rootless and rootful modes  
  Support des modes avec et sans root

**Security Scanning • Scanning de Sécurité**
- ✅ **Trivy Integration (Active):** CVE scanning for container images  
  ✅ **Intégration Trivy (Actif) :** Scan CVE pour images de conteneurs
- ✅ **Privileged Detection:** Identifies dangerous --privileged flag  
  ✅ **Détection Privilégié :** Identifie le flag --privileged dangereux
- ✅ **Sensitive Mounts:** Scans for risky filesystem bindings  
  ✅ **Montages Sensibles :** Scan des montages système risqués
- ✅ **Secret Detection:** Heuristic analysis of environment variables  
  ✅ **Détection Secrets :** Analyse heuristique des variables d'environnement
- ✅ **Network Checks:** Validates network mode configurations  
  ✅ **Vérifications Réseau :** Valide les configurations réseau

**User Interface • Interface Utilisateur**
- ✅ **Interactive TUI (Active):** Real-time dashboard with Bubble Tea  
  ✅ **TUI Interactif (Actif) :** Tableau de bord temps réel avec Bubble Tea
- **CLI Commands:** Cobra-based command system  
  **Commandes CLI :** Système de commandes basé sur Cobra

**Reporting • Génération de Rapports**
- **HTML5 Reports:** Professional reports with Risk Score calculation  
  **Rapports HTML5 :** Rapports professionnels avec calcul du Score de Risque
- **JSON Export:** Machine-readable output for CI/CD  
  **Export JSON :** Sortie lisible par machine pour CI/CD
- **Text Output:** Terminal-friendly report format  
  **Sortie Texte :** Format de rapport adapté au terminal

**Performance • Performance**
- Multi-threaded scanning with Goroutines  
  Scanning multi-threadé avec Goroutines
- Average scan time: < 10s for typical environments  
  Temps de scan moyen : < 10s pour environnements typiques

**Configuration • Configuration**
- Viper support for YAML config files  
  Support Viper pour fichiers de config YAML
- Environment variable overrides  
  Surcharge par variables d'environnement

#### 🔧 Technical Details • Détails Techniques

- **Language:** Go 1.22+ • **Langage :** Go 1.22+
- **Build:** Static binary (single file) • **Build :** Binaire statique (fichier unique)
- **Size:** ~6MB optimized • **Taille :** ~6MB optimisé
- **Platforms:** Linux, macOS, WSL2 • **Plateformes :** Linux, macOS, WSL2
- **Architecture:** amd64, arm64

---

## [Unreleased • Non publié]

### Planned for v1.1 • Prévu pour v1.1

- [ ] Auto-remediation capabilities  
      Capacités d'auto-remédiation
- [ ] Live watch mode (daemon)  
      Mode surveillance live (daemon)
- [ ] Remote scanning via SSH  
      Scanning distant via SSH
- [ ] Extended CIS Benchmark checks  
      Vérifications CIS Benchmark étendues

---

## Release Strategy • Stratégie de Release

- **v1.0.0:** Gold Master (Current • Actuel)
- **v1.x.x:** Feature updates • Mises à jour fonctionnelles
- **v2.0.0:** Major architectural changes • Changements architecturaux majeurs