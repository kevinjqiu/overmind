FROM scratch
ADD ./overmind /overmind
ENV OVERMIND_HTTP_ADDR 0.0.0.0:8080
CMD ["/overmind"]
