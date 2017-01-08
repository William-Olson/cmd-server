var express = require('express');
var router = express.Router();
const VersionManager = require('./lib/VersionManager');
const vm = new VersionManager();

const slackErr = {
  Message: 'Uh oh: Bad Request =('
};

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
router.get('/version', async (req, res, next) => {
  try {
    let url = req.query.q || null;
    const notif = await vm.createNotification(url)
    res.send(notif);
  }
  catch (err) {
    next(slackErr);
  }
});

/**
 * get slug version
 */
router.post('/slug-version', async (req, res, next) => {
  try {
    const cmd = req.body.command;
    let slug = req.body.text;

    if (!cmd.includes('/version')) {
      throw new Error('Bad Command');
    }

    if (slug.includes(' ')) {
      slug = slug.trim().split(' ')[0];
    }

    const notif = await vm.createSlugNotif(slug);
    res.send(notif);
  }
  catch (err) {
    next(slackErr);
  }
});

/*
  Fixed commands / rapid dev testing route
*/
router.post('/cmd', function(req, res, next) {
  try {
    const cmd = req.body.command;
    const params = req.body.text.split(' ');

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

