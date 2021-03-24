FROM registry.access.redhat.com/ubi8:latest
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux

RUN dnf install golang git -y

WORKDIR /go/src/github.com/mvazquezc/k8s-mutate-webhook
ADD . .
# Get dependencies
RUN go get -v ./... 
RUN go build -v -o mutateme cmd/main.go
EXPOSE 8443
CMD ["./mutateme"]
