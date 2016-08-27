FROM node:6.4

WORKDIR /opt/server

EXPOSE 7447

RUN mkdir -p /opt/server

COPY ./package.json /opt/server/package.json

RUN npm install

COPY . /opt/server

CMD ["npm", "start"]
