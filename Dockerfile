FROM golang:1.18-alpine3.16 AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ARG GITHUB_TOKEN
ENV GOPRIVATE="github.com/kargotech/*"
RUN git config --global url."https://x-access-token:$GITHUB_TOKEN@github.com".insteadOf "https://github.com"

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.* ./

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/app ./cmd/app-http/.

FROM alpine:3.16

# Make user as default user, avoiding root user as default user
ARG USERNAME=defaultuser
ARG USER_UID=1000
ARG USER_GID=$USER_UID

WORKDIR /app

# Create the user
RUN apk add sudo libcap curl bash --no-cache\
    && addgroup --gid $USER_GID $USERNAME \
    && adduser -u $USER_UID -G $USERNAME -h /home/$USERNAME -D $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    && curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 \
    && chmod +x /usr/local/bin/dbmate

# ********************************************************
# * Anything else you want to do like clean up goes here *
# ********************************************************

# [Optional] Set the default user. Omit if you want to keep the default as root.


ARG COMMIT_HASH=local

# get the hash of latest commit so we can fast and
# accurately identify the running code in deployment
RUN echo "${COMMIT_HASH}" > ./hash

COPY --from=builder /app/files ./files
COPY --from=builder /app/out/app ./app
COPY --from=builder /app/db ./db

# Enable app bind port 80 when used by non-root user
RUN setcap 'cap_net_bind_service=+ep' /app/app

USER $USERNAME

EXPOSE 80
ENTRYPOINT [ "./files/deployment/entrypoint.sh" ]
CMD [ "/app/app" ]

