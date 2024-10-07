# FROM alpine:latest
FROM archlinux:base-20240929.0.266368 AS builder

WORKDIR /app
RUN pacman -Syu --noconfirm && pacman -S go gcc --noconfirm

COPY . ./
RUN ./compile.sh

#
# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#

FROM archlinux:base-20240929.0.266368
WORKDIR /app
# RUN apk add go
# COPY kademlia ./kademlia
# COPY main.go ./main.go
# COPY go.mod ./go.mod
# ENTRYPOINT go run main.go
# COPY --chmod=0755 d7024e ./d7024e
# COPY d7024e ./d7024e
COPY --from=builder /app/d7024e /app/cli/cli ./
# COPY cli/cli ./cli
# COPY . ./
# RUN go build
ENTRYPOINT ./d7024e
# $ docker build . -t kadlab
