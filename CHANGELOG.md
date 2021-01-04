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
