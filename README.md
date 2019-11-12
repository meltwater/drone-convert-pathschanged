A conversion extension to include/exclude pipelines and steps based on paths changed.

_Please note this project requires Drone server version 1.4 or higher._

## Installation

1. Create a github token via https://github.com/settings/tokens with the scope of`repo:status`.

2. Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

3. Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=GITHUB_TOKEN=9e6eij3ckzvpe9mrhnqcis6zf8dhopmm46e3pi96 \
  --restart=always \
  --name=converter meltwater/drone-convert-pathschanged
```

4. Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_CONVERT_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
```

## Examples

This extension uses [doublestar](https://github.com/bmatcuk/doublestar) for matching paths changed in your commit range, refer to their documentation for all supported patterns.

Only run a pipeline when `README.md` is changed:
```yaml
---
kind: pipeline
name: readme

trigger:
  paths:
    include:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed”
```

Only run a pipeline step when `README.md` is changed:
```yaml
---
kind: pipeline
name: readme

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed”
  when:
    paths:
      include:
      - README.md
```

Run a pipeline step when `.yml` files are changed in the root, except for `.drone.yml`:
```yaml
---
kind: pipeline
name: yaml

steps:
- name: message
  image: busybox
  commands:
  - echo "A .yml file in the root of the repo other than .drone.yml was changed"
  when:
    paths:
      include:
      - "*.yml"
      exclude:
      - .drone.yml
```
