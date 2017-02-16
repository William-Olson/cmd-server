const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const Harness = require('route-harness');
const express = require('express');
const http = require('http');
const debug = require('debug')('server:server');

const customWrapper = require('./harness-wrapper.js');
const routes = require('./Routes.js');

const PORT = '7447';

module.exports = class Server {
  constructor()
  {
    this._listener = null;
    this._port = this._normalizePort(process.env.PORT || PORT);
    this._app = express();
    this._harness = new Harness(this._app, {
      customWrapper
    });
  }

  start()
  {
    this._init();
    this._listen();
  }

  _listen()
  {
    this._app.set('port', this._port);

    // create the HTTP server
    this._listener = http.createServer(this._app);

    // listen on provided port, on all network interfaces.
    this._listener.listen(this._port);
    this._listener.on('error', this._onListenError.bind(this));
    this._listener.on('listening', this._onListening.bind(this));
  }

  /**
   * configures express app and inits routes
   */
  _init()
  {
    // init middleware utils
    this._app.use(bodyParser.json());
    this._app.use(bodyParser.urlencoded({ extended: false }));
    this._app.use(cookieParser());

    // init routes
    this._harness.use('/', routes);

    this._initErrorHandlers();
  }

  /**
   * sets up the express app error handlers
   */
  _initErrorHandlers()
  {
    // catch 404 and forward to error handler
    this._app.use((req, res, next) => {
      const err = new Error('Not Found');
      err.status = 404;
      next(err);
    });

    // development error handler
    // will print stacktrace
    if (this._app.get('env') === 'development') {
      this._app.use((err, req, res, next) => {
        res.status(err.status || 500);
        res.send({
          message: err.message,
          error: err
        });
      });
    }

    // production error handler
    // no stacktraces leaked to user
    this._app.use((err, req, res, next) => {
      res.status(err.status || 500);
      res.send({
        message: err.message,
        error: {}
      });
    });

  }

  _onListening()
  {
    this._addr = this._listener.address();

    const bind = typeof this._addr === 'string'
      ? `pipe ${this._addr}`
      : `port ${this._addr.port}`;

    debug(`Listening on ${bind}`);
  }

  _onListenError(error)
  {
    if (error.syscall !== 'listen') {
      throw error;
    }

    const bind = typeof this._port === 'string'
      ? 'Pipe ' + this._port
      : 'Port ' + this._port;

    // handle specific listen errors with friendly messages
    switch (error.code) {
      case 'EACCES':
        debug(`${bind} requires elevated privileges`);
        process.exit(1);
        break;
      case 'EADDRINUSE':
        debug(`${bind} is already in use`);
        process.exit(1);
        break;
      default:
        throw error;
    }
  }

  _normalizePort(val)
  {
    const port = parseInt(val, 10);
    return isNaN(port) ? val : (port >= 0 ? port : false);
  }

};