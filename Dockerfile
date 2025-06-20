FROM golang:latest as api-builder

ARG ENABLE_PROXY=false

# Set working directory inside the container
WORKDIR /usr/local/apisix-dashboard

# Copy local apisix-dashboard source code into the container
COPY . ./

# Set Go proxy if needed, enable Go modules, and build the API
RUN if [ "$ENABLE_PROXY" = "true" ] ; then go env -w GOPROXY=https://goproxy.io,direct ; fi \
    && go env -w GO111MODULE=on \
    && CGO_ENABLED=0 ./api/build.sh


FROM node:16-alpine as fe-builder

ARG ENABLE_PROXY=false

WORKDIR /usr/local/apisix-dashboard

# Copy source again for frontend build
COPY . ./

WORKDIR /usr/local/apisix-dashboard/web

# Set Yarn registry if needed and build the frontend
RUN if [ "$ENABLE_PROXY" = "true" ] ; then yarn config set registry https://registry.npmmirror.com/ ; fi \
    && yarn install \
    && yarn build


FROM alpine:latest as prod

ARG ENABLE_PROXY=false

# Set Alpine mirror if proxy is enabled
RUN if [ "$ENABLE_PROXY" = "true" ] ; then sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories ; fi

WORKDIR /usr/local/apisix-dashboard

# Copy build outputs from previous stages
COPY --from=api-builder /usr/local/apisix-dashboard/output/ ./
COPY --from=fe-builder /usr/local/apisix-dashboard/output/ ./

# Create logs directory and expose the dashboard port
RUN mkdir logs

EXPOSE 9000

CMD [ "/usr/local/apisix-dashboard/manager-api" ]
