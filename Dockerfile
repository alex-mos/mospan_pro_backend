FROM golang:1.11.4-alpine3.8

RUN apk add --no-cache git g++ make sqlite
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/alex-mos/mospan_pro_backend
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go build -o mospan_pro_backend; cp mospan_pro_backend /app/
# copy database:
RUN cd $SRC_DIR; cp database.db /app/
# cleanup
RUN apk del git g++ make

ENTRYPOINT ["./mospan_pro_backend"]
