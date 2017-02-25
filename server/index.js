// use bluebird by default (clobbers global Promise)
require('any-promise/register')('bluebird', {
  Promise: require('bluebird'),
  global: true
});

const Server = require('./Server');

/**
 * Start the server
 */
const server = new Server();
server.start();
