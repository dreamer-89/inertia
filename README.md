<p align="center">
  <img src="/.static/inertia-with-name.png" width="30%"/>
</p>

<p align="center">
  Simple, self-hosted continuous deployment.
</p>

<p align="center">
  <a href="https://travis-ci.org/ubclaunchpad/inertia">
    <img src="https://travis-ci.org/ubclaunchpad/inertia.svg?branch=master"
      alt="Built Status" />
  </a>

  <a href="https://goreportcard.com/report/github.com/ubclaunchpad/inertia">
    <img src="https://goreportcard.com/badge/github.com/ubclaunchpad/inertia" alt="Clean code" />
  </a>

  <a href="https://github.com/ubclaunchpad/inertia/releases">
    <img src="https://img.shields.io/github/release/ubclaunchpad/inertia.svg" />
  </a>
</p>

----------------

Inertia is a cross-platform command line tool that simplifies setup and management of automated project deployment on any virtual private server. It aims to provide the ease and flexibility of services like Heroku without the complexity of Kubernetes while still giving users full control over their projects.

- [Usage](#rocket-usage)
  - [Setup](#setup)
  - [Continuous Deployment](#continuous-deployment)
  - [Deployment Management](#deployment-management)
  - [Release Streams](#release-streams)
- [Motivation and Design](#bulb-motivation-and-design)
- [Development](#construction-development)

# :rocket: Usage

All you need to get started is a docker-compose project, an Inertia CLI binary, and access to a virtual private server.

Inertia CLI binaries can be downloaded for various platforms in the [Releases](https://github.com/ubclaunchpad/inertia/releases) page. You can add this binary to your PATH or execute it directly to use the Inertia CLI:

```bash
$> mv $INERTIA_IMAGE ./inertia
$> ./inertia
```

You can also [install Inertia from source](#installing-from-source).

## Setup

Initializing a project for use with Inertia only takes a few simple steps:

```bash
$> inertia init
$> inertia remote add $VPS_NAME
```

After adding a remote, you can bring the Inertia daemon online on your VPS:

```bash
$> inertia $VPS_NAME init
$> inertia $VPS_NAME status
# Confirms that the daemon is online and accepting requests
```

An Inertia daemon is now running on your remote instance. This daemon will be used to manage your deployment.

See our [wiki](https://github.com/ubclaunchpad/inertia/wiki/VPS-Compatibility) for more details on platform compatibility.

## Continuous Deployment

You can now set up continuous deployment using the output of `inertia $VPS_NAME init`:

```bash
GitHub Deploy Key (add here https://www.github.com/<your_repo>/settings/keys/new):
ssh-rsa <...>
```

The Inertia daemon requires readonly access to your GitHub repository. Add the deploy key to your GitHub repository settings at the URL provided in the output - this will grant the daemon access to clone your repository.

```bash
GitHub WebHook URL (add here https://www.github.com/<your_repo>/settings/hooks/new):
http://myhost.com:8081
Github WebHook Secret: inertia
``` 

The daemon will accept POST requests from GitHub at the URL provided. Add this webhook URL in your GitHub settings area (at the URL provided) so that the daemon will receive updates from GitHub when your repository is updated.

## Deployment Management

To manually deploy your project:

```bash
$> inertia $VPS_NAME up --stream
```

There are a variety of other commands available for managing your project deployment. See the CLI documentation for more details:

```bash
$> inertia $VPS_NAME --help
```

## Release Streams

The version of Inertia you are using can be seen in Inertia's `.inertia.toml` configuration file, or by running `inertia --version`.

You can manually change the daemon version pulled by editing the Inertia configuration file. If you are building from source, you can also check out the desired version and run `make inertia-tagged`.

- `v0.x.x` denotes [official, tagged releases](https://github.com/ubclaunchpad/inertia/releases) - these are recommended.
- `latest` denotes the newest builds on `master`.
- `canary` denotes experimental builds used for testing and development - do not use this.

The daemon component of an Inertia release can be patched separately from the CLI component - see our [wiki](https://github.com/ubclaunchpad/inertia/wiki/Daemon-Releases) for more details.

## Swag

[![Deployed with Inertia](https://img.shields.io/badge/Deploying%20with-Inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)


```
[![Deployed with Inertia](https://img.shields.io/badge/Deploying%20with-Inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)
```

# :bulb: Motivation and Design

At [UBC Launch Pad](http://www.ubclaunchpad.com), we are frequently changing hosting providers based on available funding and sponsorship. Inertia is a project to develop an in-house continuous deployment system to make deploying applications simple and painless, regardless of the hosting provider.

The primary design goals of Inertia are to:

* minimize setup time for new projects
* maximimise compatibility across different client and VPS platforms
* offer a convenient interface for managing the deployed application

## How It Works

Inertia consists of two major components: a deployment daemon and a command line interface.

The deployment daemon runs persistently in the background on the server, receiving webhook events from GitHub whenever new commits are pushed. The CLI provides an interface to adjust settings and manage the deployment - this is done through requests to the daemon, authenticated using JSON web tokens generated by the daemon. Remote configuration is stored locally in `.inertia.toml`.

<p align="center">
<img src="https://bobheadxi.github.io/assets/images/posts/inertia-diagram.png" width="70%" />
</p>

Inertia is set up serverside by executing a script over SSH that installs Docker and starts an Inertia daemon image with [access to the host Docker socket](https://bobheadxi.github.io/dockerception/#docker-in-docker). This Docker-in-Docker configuration gives the daemon the ability to start up other containers *alongside* it, rather than *within* it, as required. Once the daemon is set up, we avoid using further SSH commands and execute Docker commands through Docker's Golang API. Instead of installing the docker-compose toolset, we [use a docker-compose image](https://bobheadxi.github.io/dockerception/#docker-compose-in-docker) to build and deploy user projects.

# :construction: Development

This section outlines the various tools available to help you get started developing Inertia. Run `make commands` to list all the Makefile shortcuts available.

If you would like to contribute, feel free to comment on an issue or make one and open up a pull request!

## Installing from Source

First, [install Go](https://golang.org/doc/install#install) and grab Inertia's source code:

```bash
$> go get -u github.com/ubclaunchpad/inertia
```

We use [dep](https://github.com/golang/dep) for managing Golang dependencies, and [npm](https://www.npmjs.com) to manage dependencies for Inertia's React web app. Make sure both are installed before running the following commands.

```bash
$> dep ensure     # Golang dependencies
$> make web-deps  # Web app dependencies
```

For usage, it is highly recommended that you use a [tagged build](https://github.com/ubclaunchpad/inertia/releases) to ensure compatibility between the CLI and your Inertia daemon.

```bash
$> git checkout $VERSION
$> make inertia-tagged    # installs a tagged version of Inertia
$> inertia --version      # check what version you have installed
```

Alternatively, you can manually edit `.inertia.toml` to use your desired daemon version - see the [Release Streams](#release-streams) documentation for more details.

For development, it is recommended that you make install a build tagged as `test` so that you can make use `make testdaemon` for local development. See the next section for more details.

```bash
$> make RELEASE=test
```

## Testing

You will need Docker installed and running to run the Inertia test suite, which includes a number of integration tests.

```bash
$> make test                              # test against ubuntu:latest
$> make test VPS_OS=ubuntu VERSION=14.04  # test against ubuntu:14.04
```

You can also manually start a container that sets up a mock VPS for testing:

```bash
$> make RELEASE=test
# installs current Inertia build and mark as "test"
$> make testenv VPS_OS=ubuntu VERSION=16.04
# defaults to ubuntu:lastest without args
# note the location of the key that is printed
```

You can [SSH into this testvps container](https://bobheadxi.github.io/dockerception/#ssh-services-in-docker) and otherwise treat it just as you would treat a real VPS:

```bash
$> cd /path/to/my/dockercompose/project
$> inertia init
$> inertia remote add local
# PEM file: inertia/test_env/test_key, User: 'root', Address: 0.0.0.0
$> inertia local init
$> inertia local status
```

The above steps will pull a daemon image from Docker Hub based on the version in your `.inertia.toml`.

### Inertia Daemon

To use a daemon compiled from source, set your Inertia version in `.inertia.toml` to `test` and run:

```bash
$> make testdaemon
$> inertia local init
```

This will build a daemon image and `scp` it over to the test VPS. Setting up a `testvps` using `inertia local init` will now use this custom daemon.

If you run into this error when deploying onto the `testvps`:

```bash
docker: Error response from daemon: error creating aufs mount to /var/lib/docker/aufs/mnt/fed036790dfcc73da5f7c74a7264e617a2889ccf06f61dc4d426cf606de2f374-init: invalid argument.
```

You probably need to go into your Docker settings and add this line to the Docker daemon configuration file:

```js
{
  ...
  "storage-driver" : "aufs"
}
```

This sneaky configuration file can be found under `Docker -> Preferences -> Daemon -> Advanced -> Edit File`.

### Inertia Web

To run a local instance of Inertia Web:

```bash
$> make web-run
```

Make sure you have a local daemon set up for this web app to work - see the previous section for more details.

## Compiling Bash Scripts

To bootstrap servers, some bash scripting is often involved, but we'd like to avoid shipping bash scripts with our go binary. So we use [go-bindata](https://github.com/jteeuwen/go-bindata) to compile shell scripts into our go executables.

```bash
$> go get -u github.com/jteeuwen/go-bindata/...
```

If you make changes to the bootstrapping shell scripts in `client/bootstrap/`, convert them to `Assets` by running:

```bash
$> make bootstrap
```

Then use your asset!

```go
shellScriptData, err := Asset("cmd/bootstrap/myshellscript.sh")
if err != nil {
  log.Fatal("No asset with that name")
}

// Optionally run shell script over SSH.
result, _ := remote.RunSSHCommand(string(shellScriptData))
```
