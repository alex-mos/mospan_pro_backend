FROM golang:1.11.4-alpine3.8

WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/alex-mos/mospan_pro_backend
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go build -o mospan_pro_backend; cp mospan_pro_backend /app/

ENTRYPOINT ["./mospan_pro_backend"]
