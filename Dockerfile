# Copyright 2023 FMNX Linux team.
# This code is covered by GPL license, which can be found in LICENSE file.
# Additional information could be found on official web page: https://fmnx.io/
# Email: help@fmnx.io
FROM fmnx.io/core/pack:latest AS flutter-build

RUN pack i fmnx.io/pkg/flutter
RUN sudo chmod a+rwx /opt/flutter

WORKDIR /src
COPY pubspec.yaml /src
RUN flutter pub get

COPY . /src
RUN flutter clean
RUN flutter build web

FROM golang:latest AS go-build

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src
RUN go build -o repo ./main.go

FROM fmnx.io/core/pack:latest

LABEL maintainer="dancheg97 <dangdancheg@gmail.com>"
LABEL source="https://fmnx.io/core/pack"

RUN sudo mkdir /var/cache/pacman/initpkg
RUN sudo mv -v /var/cache/pacman/pkg/* /var/cache/pacman/initpkg

COPY --from=go-build /src/repo .
COPY --from=flutter-build /src/build/web /web
RUN sudo chmod a+rwx -R /web

ENTRYPOINT ["./repo"]
CMD ["run"]