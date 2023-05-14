# 2023 FMNX team.
# Use of this code is governed by GNU General Public License.
# Additional information can be found on official web page: https://fmnx.su/
# Contact email: help@fmnx.su

FROM fmnx.su/core/pack:latest AS build

RUN pack i go fmnx.su/pkg/flutter
RUN sudo chmod a+rwx -R /opt/flutter
RUN git config --global --add safe.directory /opt/flutter

WORKDIR /home/pack
COPY pubspec.yaml /home/pack
COPY pubspec.lock /home/pack
RUN sudo flutter pub get

COPY go.mod /home/pack/
COPY go.sum /home/pack/
RUN go mod download

COPY . /home/pack

RUN sudo flutter clean
RUN sudo flutter build web
RUN go build -o repo .

FROM fmnx.su/core/pack:latest

LABEL maintainer="dancheg <dancheg@fmnx.su>"
LABEL source="https://fmnx.su/core/repo"

RUN sudo mkdir /var/cache/pacman/initpkg
RUN sudo mv -v /var/cache/pacman/pkg/* /var/cache/pacman/initpkg

COPY --from=build /home/pack/repo .
COPY --from=build /home/pack/build/web /web
RUN sudo chmod a+rwx -R /web

ENTRYPOINT ["./repo"]
CMD ["run"]