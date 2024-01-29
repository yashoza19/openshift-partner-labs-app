# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM registry.access.redhat.com/ubi9/go-toolset as builder

ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /opt/app-root/src/openshift-partner-labs-app && git config --global --add safe.directory /opt/app-root/src/openshift-partner-labs-app
WORKDIR /opt/app-root/src/openshift-partner-labs-app

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ADD . .
RUN go build -v -tags production -o /opt/app-root/src/app ./cmd/app

FROM registry.access.redhat.com/ubi9/ubi-micro:latest

WORKDIR /opt/app-root/src/

COPY --from=builder /opt/app-root/src/app .

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD exec /opt/app-root/src/app
