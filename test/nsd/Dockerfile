FROM debian:bookworm

RUN apt-get update && \
    apt-get install -y \
        autoconf \
        bison \
        flex \
        g++ \
        gcc \
        gzip \
        libevent-dev \
        libssl-dev \
        make \
        tar \
        &&  \
    apt-get clean

ADD https://github.com/NLnetLabs/nsd/archive/refs/tags/NSD_4_11_0_REL.tar.gz /app/nsd.tar.gz
ADD https://github.com/NLnetLabs/simdzone/archive/83a2327a524ffaf35c81d03e8f48073c7991944e.tar.gz /app/simdzone.tar.gz

WORKDIR /app
RUN useradd --system --user-group --home-dir /etc/nsd nsd && \
    mkdir -p /var/db/nsd && \
    chown nsd:nsd /var/db/nsd && \
    tar --strip-components=1 -xf /app/nsd.tar.gz && \
    (cd simdzone && tar --strip-components=1 -xf /app/simdzone.tar.gz) && \
    autoreconf -fi && \
    ./configure && \
    make && \
    make install

ENTRYPOINT ["/usr/local/sbin/nsd", "-d"]