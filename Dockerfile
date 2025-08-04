FROM alpine:3.19
ARG APP_NAME=VAR1
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app/

# Copy the binary file from previous stage
COPY build/${APP_NAME} /app/${APP_NAME}

# add group user
RUN addgroup --gid 1001 -S "$APP_NAME" && \
    adduser -G "$APP_NAME" --shell /bin/false --disabled-password -H --uid 1001 "$APP_NAME" && \
    mkdir -p "/var/log/$APP_NAME" && \
    chown "$APP_NAME:$APP_NAME" "/var/log/$APP_NAME"

# Copy the environment variable

# Expose application port
EXPOSE $PORT

# ser user naming
USER $APP_NAME

# Command to start the application
CMD ["/app/${APP_NAME}"]