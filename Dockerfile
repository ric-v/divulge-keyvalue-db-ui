# build image from busybox (1.5MB with shell) or alpine (5MB with full os) or scratch (0K no shell access)
FROM busybox

# copy files in the host that is to be copied to container
ADD ui/build /ui/build
COPY divulge /

# expose ports from container to network (container:host)
EXPOSE 8080:8080

# main executable for this container service
ENTRYPOINT [ "/divulge" ]
