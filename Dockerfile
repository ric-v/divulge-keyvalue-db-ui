# build image from busybox (1.5MB with shell) or alpine (5MB with full os) or scratch (0K no shell access)
FROM scratch

# copy files in the host that is to be copied to container
COPY public /public
COPY divulge-viewer-v.0.2.0-beta-amd64 /divulge-viewer-v.0.2.0-beta-amd64

# expose ports from container to network (container:host)
EXPOSE 8080:8080

# main executable for this container service
ENTRYPOINT [ "/divulge-viewer-v.0.2.0-beta-amd64" ]
