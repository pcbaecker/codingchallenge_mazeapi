# Codingchallenge MazeAPI

This repository contains an API that generates, stores and solves very simple mazes. It was built for a codingchallenge and is not intended for real world use.

### Build the docker container

On a MacBook with Apple Silicion you can use the following command to build the docker container:

    docker buildx build --platform linux/amd64 --push -t pcbaecker/codingchallenge_mazeapi:v1 .