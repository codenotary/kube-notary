# CHANGELOG
All notable changes to this project will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
<a name="unreleased"></a>
## [Unreleased]

### Bug Fixes
- handle nil option
- errors message for results endpoint
- correct image id resolution
- enable prometheus scraping
- skipping not running containers
- **dashboard:** bar chart datasource

### Changes
- verification code improvements
- update deps
- remove container_id from metrics
- use custom cluster role
- go mod tidy
- cleanup manual installation files
- basic health check
- switch to vcn exporter default port 9581
- initial support for image ID formats
- initial commit
- **grafana:** unuseful columns hidden

### Code Refactoring
- image resolution cache per container
- rename to kube-notary

### Features
- script to bulk sign all images
- config options for trusted signers, by keys or by org
- verification result endpoint
- API to get artifact from platform
- status bar chart for grafana dashboard
- grafana dashboard
- image resolutions endpoint
- image resolution caching
- cluster-wide and namespaced installations
- support for auth keychain and private registries
- manual installation templates
- helm chart
- overview dashboard for grafana
- prometheus metrics
- configmap with hot-reloading
- basic k8s deployment
- build system
- watcher
- vcn integration


[Unreleased]: https://github.com/vchain-us/kube-notary/compare/v0.1.0...HEAD
