**docker-utils**

Small utility on top of docker utility, allows you to:
- Bulk remove docker images by tag or repository name, supports wildcard matches
- Pull/push image between two docker respositories

**Installing docker-utils**
1. Using homebrew
    ```
    brew install bharat-p/tap/docker-utils
    ```

2. Download pre-built binary from:           

    https://github.com/bharat-p/docker-utils/releases


3. build your own binary
    ```
   go get github.com/bharat-p/docker-utils

   cd $GOPATH/src/github.com/bharat-p/docker-utils

   go build -o docker-utils main.go 
   ```

**Usage**

    docker-utils -h
    Provide some utility commands to use with docker.

    Usage:
    docker-utils [command]

    Available Commands:
    help        Help about any command
    pull-push   Pull and push images between repositories
    rmi         Remove docker images by specifying repository name or tag name
    version     Print the version

    Flags:
    -D, --dry-run   Dry Run
    -V, --verbose   Debug mode
    
    
***docker-utils rmi***

Bulk delete docker images from local system

```
docker-utils rmi -h  
Remove docker images by specifying repository name or tag name, For example
:
To remove all images with repository ubuntu
docker-utils rmi -r "ubuntu"

Usage:
  docker-utils rmi [flags]

Flags:
  -f, --force               Force removal of image(s)
  -r, --repository string   Repository name to remove images for
  -t, --tag string          Remove all image matching tag

Global Flags:
  -D, --dry-run   Dry Run
  -V, --verbose   Debug mode

```

Examples:

* Remove all images with tag v1.0

    ```
    docker-utils rmi -t "v1.0"
    ```

    _use flag -f to force removal of images_

* Remove all images with repo name: `release/myapp` 

    ```
    docker-utils rmi -r "release/myapp"
    ```    
* Remove all images where repo name startw with: `release/` 
    ```
    docker-utils rmi -D -r "release/.*"
    ```

   _-D is for dry run mode, won't delete any image, will only print command that will get executed_

***docker-utils pull-push***

Pull/Push multiple images from one repository to another.

```
docker-utils pull-push -h  
Pull images from one docker registry and push it to another.

Usage:
  docker-utils pull-push [flags]

Flags:
  -f, --from string          Registry from where to pull images
  -i, --images stringArray   Image(s) to pull/push
  -l, --local                Tag as local (applicable with --from only if --to is not used)
  -r, --remove               Remove images after pull/push is done.
  -t, --to string            Registry to push images

Global Flags:
  -D, --dry-run   Dry Run
  -V, --verbose   Debug mode

```

Examples:

* Pull images: `my-image1:v1.0` and `my-image2` from registry `registry.example.com` locally

    ```
    docker-utils pull-push -i my-image1:v1.0,my-image2 \ 
        --from registry.example.com
    ```

    _images will be pulled as registry.example.com/my-image:v1.0 and registry.example.com/my-image2:latest, if you wish to remove remote registry name from image name when pulling image, please specify flag -l or --local_



* Pull images: `my-image1:v1.0` and `my-image2` from registry `registry.example.com` and push to another registry `registry.local`

    ```
    docker-utils pull-push -i my-image1:v1.0,my-image2 \ 
        --from registry.example.com \
        --to registry.local
    ```

