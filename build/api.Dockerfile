FROM alpine
RUN apk add --no-cache ca-certificates
ADD api /app/api
CMD ["/app/api"]