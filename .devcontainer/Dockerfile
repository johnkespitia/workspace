FROM mcr.microsoft.com/devcontainers/go:1

# Install common tools and Node.js LTS
USER root
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    gnupg \
    build-essential \
    git \
    procps \
  && rm -rf /var/lib/apt/lists/*

# Install Node.js LTS via NodeSource
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
  && apt-get install -y nodejs \
  && npm --version \
  && corepack enable || true

USER vscode

ENV GOPATH=/home/vscode/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH
