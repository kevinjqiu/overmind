FROM scratch
ADD ./overmind /overmind
ENV HTTPADDR 0.0.0.0:8080
CMD ["/overmind", "-http.addr=$(echo -n $HTTPADDR)"]
