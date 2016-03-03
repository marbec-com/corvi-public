# Project Corvi - Social Learning Tool

## Repository Structure

### Branches
* develop (current development)
* master (tagged releases)

### Folders
* electron
  * main.js, package.json
* frontend
  * index.html, css, angular etc.
* backend
  * main.go etc.
* .gitlab-ci.yml
* README.md

## Continuous Delivery Stages

1. Clean Up
2. Test Go
3. Build Go executable (darwin x64)
4. Build Go executable (linux x64)
5. Build Go executable (windows x64)
6. Build electron package (darwin x64)
7. Build electron package (linux x64)
8. Build electron package (windows x64)
9. Zip packages and upload to S3

### Build Docker image
* golang
* gox
* node.js + npm
* electron-packager
* electron-builder
* python + pip
* aws-cli