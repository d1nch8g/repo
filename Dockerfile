# Copyright 2023 FMNX Linux team.
# This code is covered by GPL license, which can be found in LICENSE file.
# Additional information could be found on official web page: https://fmnx.io/
# Email: help@fmnx.io
FROM fmnx.io/core/pack:latest AS build

RUN pack r python libns1
RUN pack i go fmnx.io/pkg/flutter
RUN sudo chmod a+rwx -R /opt/flutter
RUN git config --global --add safe.directory /opt/flutter

WORKDIR /home/pack
COPY pubspec.yaml /home/pack
COPY pubspec.lock /home/pack
RUN flutter pub get

COPY go.mod /home/pack/
COPY go.sum /home/pack/
RUN go mod download

COPY . /home/pack

RUN flutter clean
RUN flutter build web
RUN go build -o repo ./main.go

FROM fmnx.io/core/pack:latest

LABEL maintainer="dancheg97 <help@fmnx.io>"
LABEL source="https://fmnx.io/core/repo"

RUN sudo mkdir /var/cache/pacman/initpkg
RUN sudo mv -v /var/cache/pacman/pkg/* /var/cache/pacman/initpkg

COPY --from=build /home/pack/repo .
COPY --from=build /home/pack/build/web /web
RUN sudo chmod a+rwx -R /web

ENTRYPOINT ["./repo"]
CMD ["run"]