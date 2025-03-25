FROM golang:1.18.2-alpine3.15 as build
RUN apk --no-cache add tzdata
WORKDIR /app
ADD source/. ./
RUN ls /app
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM scratch as finals
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /app/app .
ENV TZ=Asia/Jakarta
EXPOSE 8900
CMD [ "/app" ]
