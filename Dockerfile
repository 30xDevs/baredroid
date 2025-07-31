FROM golang:1.24.5 AS builder

LABEL "dev.baredroid.xochitla"="Xochitla"
LABEL version="1.0"
LABEL description="Image for baredroid"

# Install adb and aapt/aapt2
RUN apt-get update && apt install -y adb aapt

# Clone v1.0 tag for baredroid
RUN git clone https://github.com/30xDevs/baredroid.git

# Build baredroid binary
RUN cd baredroid && go build -o brdrd

# Expose adb port
EXPOSE 5037
