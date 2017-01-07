var express = require('express');
var router = express.Router();
const VersionManager = require('./lib/VersionManager');
const vm = new VersionManager();

/*
  TODO: implement some token auth
*/



/*
  Root route
*/
router.get('/', function(req, res, next) {
  try {
    res.send({ ok: true });
  }
  catch (err) {
    next(err);
  }
});



/*
  General Version command route
*/
router.get('/cmd/version', function(req, res, next) {
  try {
    let url = req.query.q || null;
    vm.createNotification(url)
      .then(notif => res.send(notif))
      .catch(err => {
        res.send({ error: 'An error occurred!', err });
      });
  }
  catch (err) {
    next(err);
  }
});

/**
 * get slug version
 */
router.post('/slug-version', async (req, res, next) => {
  try {
    let params = req.body.text;
    const cmd = req.body.command;

    if (!cmd.includes('/version')) {
      throw new Error('Bad Command');
    }

    if (params.includes(' ')) {
      params = params.trim().split(' ')[0];
    }

    const notif = await vm.createSlugNotif(params);
    res.send(notif);
  }
  catch (err) {
    next(err);
  }
});

/*
  Fixed commands / rapid dev testing route
*/
router.post('/cmd', function(req, res, next) {
  try {
    const cmd = req.body.command;
    const params = cmd.split(' ');

    if (!params.length) {
      throw new Error('Must provide params');
    }

    switch (params[0]) {
      case '/cmd1':
        // handle example command 1
        res.send({ ok: true, cmd: params[0] });
        break;
      case '/cmd2':
        // handle example command 2
        res.send({ ok: true, cmd: params[0] });
        break;
      default:
        throw new Error(`Bad param ${params[0]}`);
        break;
    }
  }
  catch (err) {
    next(err);
  }
});



module.exports = router;

