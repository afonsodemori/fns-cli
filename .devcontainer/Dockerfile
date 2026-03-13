FROM ghcr.io/afonsodemori/devcontainers/golang:1.26

RUN apt-get update \
    && apt-get install -y nfpm \
    && apt-get autoremove -y && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*
