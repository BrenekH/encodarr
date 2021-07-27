# How to contribute

Thanks for considering to contribute to Encodarr! It means a lot.

## Starting off

If you've noticed a bug or have a feature request, [file an issue](https://github.com/BrenekH/encodarr/issues/new/choose)!
It's generally best if you get confirmation of your bug or approval for your feature request before starting to work on code.

If you have a general question about Encodarr, you can post it to [GitHub Discussions](https://github.com/BrenekH/encodarr/discussions), the issue tracker is only for bugs and feature requests.

## Code Structure

The repo is broken up into 3 main folders: `controller`, `runner`, and `frontend`.

`controller` contains code relating to the Controller, `runner` is for the Runner, and `frontend` holds the React project for the Web UI that is embedded into the Controller.

## Controller

The Controller can be built locally and ran using the standard Go build commands, `go build cmd/main.go` and `go run cmd/main.go`, from within the `controller` folder.
Unit tests can be ran using `go test ./...`.

## Runner

Like the Controller, the Runner can be built locally and ran using the standard Go build commands, `go build cmd/EncodarrRunner/main.go` and `go run cmd/EncodarrRunner/main.go`, from within the `runner` folder.
Unit tests can be ran using `go test ./...`.

## Frontend

The frontend code can be run in debug mode using `npm run start`.
It is recommended to run a Controller using the default port on your local machine so that the development Web UI can make requests to it.
If you want to point React at a different Controller, you can modify the `proxy` field in the `package.json` file, just make sure to change it back before committing anything or submitting changes.

To "deploy" the frontend to the Controller, run `npm run build` and copy the generated `build` folder to `controller/user_interfacer/webfiles`.
To perform both actions with one command on Linux systems, use `npm run build && cp -r build/. ../controller/user_interfacer/webfiles/` from within the `frontend` directory.

## Docker

Both the Controller and the Runner have `Dockerfile`s in their folders that can be built into container images using `docker build .` from the appropriate directory.

## Continuous Integration and Continuous Deployment

This repository uses GitHub Actions to continuously test commits and deploy them to Docker Hub and the GitHub Container Registry.

Stand-alone builds as a result of action workflows can be found in the artifacts section of the workflow run's page.
Artifacts are kept for 90 days, after which they are removed.

When a Git tag is pushed to the repo, special actions are run that deploy the code as a release, both to the docker registries(Docker Hub and GitHub Container Registry) and the [GitHub Releases page](https://github.com/BrenekH/encodarr/releases).

## Code of Conduct

This project holds all maintainers, contributors, and participants to the standards outlined by the Contributor Covenant, a copy of which can be found in [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
