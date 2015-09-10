FROM golang

ENV GO15VENDOREXPERIMENT=1

ADD ./api /go/src/grubprint.io/api
ADD ./app /go/src/grubprint.io/app
ADD ./client /go/src/grubprint.io/client
ADD ./cmd /go/src/grubprint.io/cmd
ADD ./datastore /go/src/grubprint.io/datastore
ADD ./httputil /go/src/grubprint.io/httputil
ADD ./keystore /go/src/grubprint.io/keystore
ADD ./router /go/src/grubprint.io/router
ADD ./usda /go/src/grubprint.io/usda
ADD ./vendor /go/src/grubprint.io/vendor

RUN go install grubprint.io/cmd/grubprint
RUN /go/bin/grubprint -keygen

ENTRYPOINT /go/bin/grubprint
EXPOSE 8080
