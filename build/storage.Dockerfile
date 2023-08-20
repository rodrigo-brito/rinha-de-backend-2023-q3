FROM alpine
RUN apk add --no-cache ca-certificates
ADD storage /app/storage
CMD ["/app/storage"]