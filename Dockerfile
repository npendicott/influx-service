FROM golang:1.12-alpine
#FROM golang:1.8-alpine  # This should work, issues with arm tho

# Deps
RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/influx-client-london
COPY . .

RUN go get -v 
RUN go install -v .

# Should probably not have this in the DF, 
# instead just manage my .env better
ENV ENERGY_DB_HOST 'influx'

CMD ["influx-client-london"]
