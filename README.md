# rio-sample is a sample app to show case

The sample app generates a static web page with a different background color.

The sample helps show case the following:

- Building Apps with Rio
- Deploying a new service with Rio
- Deploying a new version of the service with Rio
- Canary release with Rio
- Switching traffic with Rio
- Rolling back with Rio

### Building Apps with Rio
In the directory containing the Dockerfile run the following
`rio build -t rio-sample:white --build-arg COLOR="white"`

We should see something similar to this:
```
▶ rio build -t rio-sample:white --build-arg COLOR=white
INFO[0005] Building service rio-sample
[+] Building 9.2s (10/10) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                                                                                                                            0.8s
 => => transferring dockerfile: 32B                                                                                                                                                                                                                                             0.8s
 => [internal] load .dockerignore                                                                                                                                                                                                                                               0.8s
 => => transferring context: 2B                                                                                                                                                                                                                                                 0.8s
 => [internal] load metadata for docker.io/library/golang:1.13                                                                                                                                                                                                                  3.1s
 => [builder 1/4] FROM docker.io/library/golang:1.13@sha256:43f859b58af8c84c8aef288809204cfbd7cb88dbd4b0cf473dd4fb86693403ad                                                                                                                                                    0.0s
 => => resolve docker.io/library/golang:1.13@sha256:43f859b58af8c84c8aef288809204cfbd7cb88dbd4b0cf473dd4fb86693403ad                                                                                                                                                            0.0s
 => [internal] load build context                                                                                                                                                                                                                                               1.3s
 => => transferring context: 1.57kB                                                                                                                                                                                                                                             1.3s
 => CACHED [builder 2/4] RUN mkdir -p /src/github.com/ibrokethecloud/rio-sample                                                                                                                                                                                                 0.0s
 => [builder 3/4] COPY . /src/github.com/ibrokethecloud/rio-sample                                                                                                                                                                                                              0.0s
 => [builder 4/4] RUN cd /src/github.com/ibrokethecloud/rio-sample     && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /root/rio-sample                                                                                                                                    3.3s
 => CACHED [stage-1 1/2] COPY --from=builder /root/rio-sample /rio-sample                                                                                                                                                                                                       0.0s
 => exporting to image                                                                                                                                                                                                                                                          0.1s
 => => exporting layers                                                                                                                                                                                                                                                         0.0s
 => => exporting manifest sha256:bfceee10ceaf2e36f9c823d84fd6f2cfd4171e51faec7abbb07c4918dc24efa8                                                                                                                                                                               0.0s
 => => exporting config sha256:ef08900938f636edc0b9f2621367c356029e7d966a902758c4165b256eab8bc5                                                                                                                                                                                 0.0s
 => => pushing layers                                                                                                                                                                                                                                                           0.0s
 => => pushing manifest for localhost:80/rio-system/rio-sample:white                                                                                                                                                                                                            0.0s
localhost:5442/rio-system/rio-sample:white
```

The list of built images can be inspected using the `rio images` command

Should list all images available in local registry

```
▶ rio images
REPO          TAG          IMAGE
hello-world   latest       localhost:5442/rio-system/hello-world:latest
rio-sample    white        localhost:5442/rio-system/rio-sample:white
rio-sample    latest       localhost:5442/rio-system/rio-sample:latest
```

### Deploying a new service with Rio
To launch a new service with Rio the following command is sufficient

`rio run -p 8080 -n rio-sample localhost:5442/rio-system/rio-sample:white`

We can now view the status of the service using `rio ps`

Once the service is ready, the endpoint for the service can be viewed using `rio endpoints`

The application can be accessed using the endpoint available for the service.

### Deploying a new version of service with Rio
We will build a new rio-sample image with a different tag and build argument for the color.

This can be accomplished as
`rio build -t rio-sample:pink --build-arg COLOR=pink`

Once the build is complete we can stage the image to the service.
`rio stage --image localhost:5442/rio-system/rio-sample:pink rio-sample pink`

Once the service is deployed, we can see the status using `rio ps`

The staged service will not have any requests sent to it, since the weight for the service is 0.

### Canary release with Rio
We can now allow a weight % of the ingress traffic to the new *pink* image we built in the last step.

This can be done as follows:
`rio weight rio-sample@pink=25`

This will set the weight for the pink colored service to 25 and will allow some traffic to be routed to the new image.

The user can verify the same by hitting the endpoint for the *rio-sample* service available from `rio endpoints` command.

### Promoting a new release with Rio
Once the user is happy with the new version of the app, we can update the weight of the service version to 100% to switch all incoming traffic to the new service.
`rio weight rio-sample@pink=100`

### Rolling back with Rio
Rolling back to an old version of the app is as easy as setting the weight of the old service to 100.
`rio weight rio-sample@white=100`

If the users are happy with the new service, the old service can be removed using *rm* command
`rio rm rio-sample@white`
