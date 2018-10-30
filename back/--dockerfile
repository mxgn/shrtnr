# STEP 1 build executable binary
FROM golang:alpine as builder

# Install git
RUN apk update && apk add git 

# Create appuser
RUN adduser -D -g '' gmaxusr

COPY . $GOPATH/src/github.com/mxgn/url-shrntr/
WORKDIR $GOPATH/src/github.com/mxgn/url-shrntr/

#get dependancies
#you can also use dep
RUN go get -d -v

#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/url-shrntr-test


# STEP 2 build a small image
# start from scratch
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/url-shrntr-test /go/bin/url-shrntr-test

USER gmaxusr

ENTRYPOINT ["/go/bin/url-shrntr-test"]

# Document that the service listens on port 8080
EXPOSE 8080