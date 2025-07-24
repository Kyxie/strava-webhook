FROM node:18-alpine

WORKDIR /app

RUN apk add --no-cache curl libc6-compat

RUN curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-24.0.7.tgz | tar xz && \
    mv docker/docker /usr/bin/docker && \
    chmod +x /usr/bin/docker && \
    rm -rf docker

COPY package*.json ./
RUN npm install

COPY . .

VOLUME /var/run/docker.sock

EXPOSE 8001

CMD ["npm", "start"]
