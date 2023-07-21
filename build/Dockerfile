FROM golang:1.19-alpine as build
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN cd cmd/hooks && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian11
COPY --from=build /main .
CMD ["./main"]
