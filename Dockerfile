FROM golang:latest AS go-build

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src
RUN go build -o repo ./main.go

FROM cirrusci/flutter AS flutter-build

WORKDIR /src
COPY pubspec.yaml /src
RUN flutter pub get

COPY . /src
RUN flutter clean
RUN flutter build web

FROM fmnx.io/core/pack:latest

COPY --from=go-build /src/repo .
COPY --from=flutter-build /src/build/web /web
RUN sudo chmod a+rwx -R /web

ENTRYPOINT ["./repo"]
CMD ["run"]