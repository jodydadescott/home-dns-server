FROM fedora:latest
WORKDIR /

RUN dnf -y update

RUN dnf install -y iproute iputils bind-utils file hostname procps net-tools dnf-plugins-core findutils

ADD build/linux-amd-64/unifi-dns-server /usr/sbin/unifi-dns-server
RUN chmod +x /usr/sbin/unifi-dns-server

# ENTRYPOINT [ "/usr/sbin/unifi-dns-server" ]
CMD ["/usr/sbin/unifi-dns-server", "run"]