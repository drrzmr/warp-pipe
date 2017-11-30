# Docker base images

The base images needed to build, develop, test and run in production environment are created here.
To build then:

```
$ docker-compose build
```

or

```
$ make
```


All images will take de prefix **warp-pipe**, so to see generated images after build:

```
$ docker images | grep ^warp-pipe
```

or

```
$ make show
```
