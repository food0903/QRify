FROM node:20-alpine

WORKDIR /app
COPY . .
RUN npm install

RUN npm run build

# Expose port 3000
EXPOSE 3000

CMD ["npm", "run", "start"]