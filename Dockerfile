FROM node:18-alpine

# Set working directory
WORKDIR /app

# Install bash (alpine only comes with ash by default)
RUN apk add --no-cache bash

# Install dependencies
COPY package.json ./
RUN npm install

# Copy all source files including update.sh
COPY . .

# Expose webhook port
EXPOSE 8001

# Start the app
CMD ["npm", "start"]
