FROM node:alpine

WORKDIR /usr/src/app

COPY . .

RUN npm install --only=production

ENV PORT 8080
EXPOSE $PORT

CMD [ "npm", "start" ]