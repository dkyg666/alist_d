#!/usr/bin/env sh

chown -R "${PUID}":"${PGID}" /opt/alist/

umask "${UMASK}"

exec su-exec "${PUID}":"${PGID}" ./alist "$@"
