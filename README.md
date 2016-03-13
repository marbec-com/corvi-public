# Project Corvi - Social Learning
Corvi is a social learning app which works similar to an index-card box.
Users can create boxes to which they can add questions like vocabulary, formulae, definitions or similar.
Those questions get tested in regular intervals, depending on the learning progress the user makes.
Corvi makes sure, that new questions are learned stepwise to keep daily learning sessions short.
At the same time it takes care that the user does not unlearn previous questions.

![Corvi](https://zippy.gfycat.com/AridFluidIberianmole.gif)

# Download
Corvi is currently available for Windows, Mac OSX and Linux.
You can find the release notes and the download of the latest version in the [release](https://github.com/marbec-com/corvi-public/releases) section.

# Repository Structure
## Branches
The main repository is located at [GitLab](https://gitlab.solid.marb.ec/marbec/corvi) and the `master` branch is also mirrored at [GitHub](https://github.com/marbec-com/corvi-public).

* `develop` (current development)
* `master` (Releases)

## Folders
This repository contains the source code of the Corvi desktop client and is divided in three folders.

Folder | Content
------------ | -------------
`backend` | The Go backend takes care of storing and synchronizing data with the frontend. It persists data using SQLite, provides a RESTful API with WebSocket notifications and delivers static files (esp. the frontend). 
`electron` | Configuration and starting script for electron. Electron takes care of starting the Go backend.
`frontend` | User interface, based on AngularJS.

## Continuous Delivery
Automated building and testing is available via GitLab CI as configured in `.gitlab-ci.yml`.
The current build environment consists of three dedicated runners - one for each platform (Linux, Mac & Windows).

1. Run Backend tests (Linux, Mac, Windows)
2. Create GitHub release
3. Build Go executable (Linux, Mac, Windows)
4. Pack Electron app (Linux, Mac, Windows)
5. Zip and upload to GitHub

## Dependencies
To build Corvi, the following dependencies are required:

* [Go](https://github.com/golang/go)
* [node.js + npm](https://github.com/nodejs/node)
* [electron-packager](https://github.com/electron-userland/electron-packager)
* [github-release](https://github.com/aktau/github-release)
* zip
* gcc (for building sqlite)
* cygwin (Windows only)

For linux, the building environment is available as Docker [container](https://github.com/herzog31/corvi-build).