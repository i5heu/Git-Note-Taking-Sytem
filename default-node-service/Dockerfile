FROM node:latest

WORKDIR /usr/src/app
COPY . .

RUN npm install
RUN npm run build

EXPOSE 80
CMD [ "npm", "start" ]