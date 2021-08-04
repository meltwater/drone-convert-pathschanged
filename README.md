[![Build Status](https://cloud.drone.io/api/badges/meltwater/drone-convert-pathschanged/status.svg)](https://cloud.drone.io/meltwater/drone-convert-pathschanged)
[![Docker Pulls](https://img.shields.io/docker/pulls/meltwater/drone-convert-pathschanged)](https://hub.docker.com/r/meltwater/drone-convert-pathschanged)

A [Drone](https://drone.io/) [conversion extension](https://docs.drone.io/extensions/conversion/) to include/exclude pipelines and steps based on paths changed.

_Please note this project requires Drone server version 1.4 or higher._

## Installation

## Github
1. Create a github token via https://github.com/settings/tokens with the scope of`repo` (see [issue 13](https://github.com/meltwater/drone-convert-pathschanged/issues/13) for background).

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
  --env=TOKEN=9e6eij3ckzvpe9mrhnqcis6zf8dhopmm46e3pi96 \
  --env=PROVIDER=github \
  --restart=always \
  --name=converter meltwater/drone-convert-pathschanged
```

4. Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_CONVERT_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
```

## Bitbucket Server

_Bitbucket Server support is currently considered experimental_

1. Create a BitBucket access token via https://your-bitbucket-address/plugins/servlet/access-tokens/manage with read-only rights

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
  --env=TOKEN=9e6eij3ckzvpe9mrhnqcis6zf8dhopmm46e3pi96 \
  --env=PROVIDER=bitbucket-server \
  --env=BB_ADDRESS=https://your-bitbucket-address \
  --restart=always \
  --name=converter meltwater/drone-convert-pathschanged
```

4. Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_CONVERT_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
```

## Gitea

_Gitea support is currently considered experimental_

1. Create a Gitea access token via https://your-bitbucket-address/plugins/servlet/access-tokens/manage with read-only rights

2. Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```


3. Download and run the plugin (using self-hosted GitLab)
```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=TOKEN=9e6eij3ckzvpe9mrhnqcis6zf8dhopmm46e3pi96 \
  --env=PROVIDER=gitea \
  --env=GITEA_ADDRESS=https://gitea.example.com \
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

### `include`

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

Same as above, but with an implicit `include`:
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
    - README.md
```

### `include` and `exclude`

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

### `depends_on`

When using [`depends_on`](https://docker-runner.docs.drone.io/configuration/parallelism/) in a pipeline step, ensure the `paths` rules match, otherwise your steps may run out of order.

Only run two steps when `README.md` is changed, one after the other:
```yaml
---
kind: pipeline
name: depends_on

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed”
  when:
    paths:
      include:
      - README.md

- name: depends_on_message
  depends_on:
  - message
  image: busybox
  commands:
  - echo "This step runs after the message step"
  when:
    paths:
      include:
      - README.md
```

## Known issues

### YAML anchors

There is a problem in the YAML library where ordering matters during unmarshaling, see https://github.com/meltwater/drone-convert-pathschanged/issues/18

This syntax will fail:

```yaml
anchor: &anchor
  image: busybox
  settings:
    foo: bar

  - name: test
    <<: *anchor
    when:
      event: push
      branch: master
```

But this will succeed:

```yaml
anchor: &anchor
  image: busybox
  settings:
    foo: bar

  - <<: *anchor
    name: test
    when:
      event: push
      branch: master
```
