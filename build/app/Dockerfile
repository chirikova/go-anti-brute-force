FROM golang:1.21 as builder
ENV NAME "app"

WORKDIR /out/${NAME}
COPY go.mod .
COPY go.sum .
RUN go mod download


FROM builder AS build
ENV NAME "app"
WORKDIR /out/${NAME}
COPY . .
RUN CGO_ENABLED=0 \
    go build \
    -o bin/${NAME} cmd/app/main.go


FROM alpine
ARG CONFIG_FILE
ENV CONFIG_FILE=$CONFIG_FILE
ENV NAME "app"
WORKDIR /out/${NAME}
COPY --from=build /out/${NAME}/configs ./configs
COPY --from=build /out/${NAME}/bin/${NAME} ./${NAME}
RUN apk add --no-cache tzdata
CMD ./${NAME} -config ${CONFIG_FILE}