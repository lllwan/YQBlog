FROM golang:1.15.5
RUN mkdir -p /app
COPY . /app
WORKDIR /app/
RUN export GOPROXY=https://goproxy.io,direct && go mod tidy
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o YQBlog .

FROM alpine
MAINTAINER YQBlog
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
RUN apk add --no-cache tzdata
RUN mkdir -p /app/templates
WORKDIR /app/
COPY --from=0 /app/YQBlog /app/YQBlog
COPY --from=0 /app/config.yaml /app/config.yaml
COPY --from=0 /app/themes/* /app/themes/
RUN chmod +x YQBlog
CMD /app/YQBlog
