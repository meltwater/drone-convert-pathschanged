## 1.0.0

### Breaking changes
- [#79](https://github.com/meltwater/drone-convert-pathschanged/pull/79) shift to drone/go-scm library support

### Added
- [#83](https://github.com/meltwater/drone-convert-pathschanged/pull/83) add missing environnment variable to stash (bitbucket-server) example
- [#77](https://github.com/meltwater/drone-convert-pathschanged/pull/77) add bitbucket cloud support
- [#76](https://github.com/meltwater/drone-convert-pathschanged/pull/76) golang module updates
- [#74](https://github.com/meltwater/drone-convert-pathschanged/pull/74) update drone/go and drone/go-scm modules
- [#73](https://github.com/meltwater/drone-convert-pathschanged/pull/74) golang tidy

## 0.6.0
### Added
- [#72](https://github.com/meltwater/drone-convert-pathschanged/pull/72) add GitHub Server support

## 0.5.0
### Breaking changes
- [#47](https://github.com/meltwater/drone-convert-pathschanged/pull/47) build images from 'scratch'
  - Official documentation https://docs.docker.com/develop/develop-images/baseimages/#create-a-simple-parent-image-using-scratch

## 0.4.0
### Breaking changes
- [#46](https://github.com/meltwater/drone-convert-pathschanged/pull/46) handle `--allow-empty` commits with excludes
  - Until this fix, `--allow-empty` commits with pipeline or path exclusion rules would have caused the pipeline or path to be excluded. With this fix, exclusions are handled properly, so empty commits will trigger pipelines and pipeline steps with path exclusion rules (since no files were changed, pipelines and steps should not be excluded). This could be potentially unexpected behavior if you had been relying on the exclusion of pipelines or steps when making empty commits.

## 0.3.1
### Added
- linux-arm64 build

## 0.3.0
### Added
- [#38](https://github.com/meltwater/drone-convert-pathschanged/pull/38) Make include implicit and is optional
    - This addresses [#33](https://github.com/meltwater/drone-convert-pathschanged/issues/33) to make this plugin compatble with [`drone jsonnet`](https://docs.drone.io/pipeline/scripting/jsonnet/)
- [#40](https://github.com/meltwater/drone-convert-pathschanged/pull/40) Additional tests

## 0.2.0
### Added
- [#29](https://github.com/meltwater/drone-convert-pathschanged/pull/29) Initial experimental Bitbucket Server support.
    - This includes breaking changes, environment variable `GITHUB_TOKEN` has been replaced with `TOKEN`, and `PROVIDER` is a new required environment variable which must be either `github` or `bitbucket-server`.

## 0.1.0
### Added
- Initial prometheus support
