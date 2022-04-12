FROM golang:latest
RUN mkdir -p /go/src/assignment
WORKDIR /go/src/assignment
COPY . /go/src/assignment
RUN go install assignment
CMD /go/bin/assignment
EXPOSE 8000