# Diskmon

Diskmon will notify you if a disk has reached a configurable size limit.

## Get started

### Using Binary

Build the binary or run directly using

```sh
go run main.go -basedir /home
```

### Using Docker

Build the image

```sh
docker build -t dockermon .
```

Run the image

```sh
docker run --volume /home:/home:ro diskmon -basedir /hom
```

## Limitations
