# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change.

Please note we have a [code of conduct](CODE_OF_CONDUCT.md) that we ask you to follow in all your interactions with the project.

**IMPORTANT: Please do not create a Pull Request without creating an issue first.**

*Any change needs to be discussed before proceeding. Failure to do so may result in the rejection of the pull request.*

Thank you for your pull request. Please provide a description above and review
the requirements below.

## Pull Request Process

0. Check out [Pull Request Checklist](#pull-request-checklist), ensure you have fulfilled each step.
1. Check out guidelines below, the project tries to follow these, ensure you have fulfilled them as much as possible.
    * [Effective Go](https://golang.org/doc/effective_go.html)
    * [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. Ensure any install or build dependencies are removed before the end of the layer when doing a
   build.
3. Please ensure the [README](README.md) and [DOCS](./DOCS.md) are up-to-date with details of changes to the command-line interface,
    this includes new environment variables, exposed ports, used file locations, and container parameters.
4. **PLEASE ENSURE YOU DO NOT INTRODUCE BREAKING CHANGES.**
5. **PLEASE ENSURE BUG FIXES AND NEW FEATURES INCLUDE TESTS.**
6. You may merge the Pull Request in once you have the sign-off of one other maintainer/code owner,
   or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## Pull Request Checklist

- [x] Read the **CONTRIBUTING** document. (It's checked since you are already here.)
- [ ] Read the [**CODE OF CONDUCT**](CODE_OF_CONDUCT.md) document.
- [ ] Add tests to cover changes.
- [ ] Ensure your code follows the code style of this project.
- [ ] Ensure CI and all other PR checks are green OR
    - [ ] Code compiles correctly.
    - [ ] Created tests which fail without the change (if possible).
    - [ ] All new and existing tests passed.
- [ ] Add your changes to `Unreleased` section of [CHANGELOG](CHANGELOG.md).
- [ ] Improve and update the [README](README.md) (if necessary).


## Release Process

*Only concerns maintainers/code owners*

0. **PLEASE DO NOT INTRODUCE BREAKING CHANGES**
1. Update `README.md`with the latest changes.
2. Increase the version numbers in any examples files and the README.md to the new version that this
   the release would represent. The versioning scheme we use is [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/meltwater/drone-convert-pathschanged/tags).

3. Ensure [CHANGELOG](CHANGELOG.md) is up-to-date with new version changes.
4. Update version references.
5. Create a tag on the master. This will trigger drone build and new images are pushed into DockerHub with the version references.

    ```console
    $ git tag -am 'vX.X.X'
    > ...
    $ git push --tags
    > ...
    ```

> **Keep in mind that users usually use the `latest` tagged images in their pipeline, please make sure you do not interfere with their working workflow.**

## Testing Locally

You can build and run the plugin locally through the CLI (you will need a GitHub token with the ‘repo’ scope)

```bash
$ cd drone-convert-pathschanged
$ openssl rand -hex 16
2ea1d6ca0df30bf3957ad1c0de441f0d
$ ./scripts/build.sh
$ docker build -t pathschanged -f docker/Dockerfile.linux.amd64 .
$ docker run --rm -e DRONE_DEBUG=true -e DRONE_SECRET=2ea1d6ca0df30bf3957ad1c0de441f0d -e PROVIDER=github -e TOKEN=REDACTED --name=converter -p 3000:3000 -it pathschanged
```

Then you can send to this endpoint by using plugins comandset:

```bash
$ export DRONE_CONVERT_SECRET=2ea1d6ca0df30bf3957ad1c0de441f0d
$ export DRONE_CONVERT_ENDPOINT=http://localhost:3000
$ drone plugins convert --path .drone.yml --before _SHA_ --after _SHA_ --ref refs/heads/changeset --repository meltwater/some-repo
```

## Response Times

**Please note the below timeframes are response windows we strive to meet. Please understand we may not always be able to respond in the exact timeframes outlined below**
- New issues will be reviewed and acknowledged with a message sent to the submitter within two business days
    - ***Please ensure all of your pull requests have an associated issue.***
- The ticket will then be groomed and planned as regular sprint work and an estimated timeframe of completion will be communicated to the submitter.
- Once the ticket is complete, a final message will be sent to the submitter letting them know work is complete.

***Please feel free to ping us if you have not received a response after one week***



