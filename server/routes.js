var express = require('express');
var router = express.Router();
const VersionManager = require('./lib/VersionManager');

const vm = new VersionManager();

/* GET home page. */
router.get('/', function(req, res, next) {
  try {
    res.send({ok: true});
  }
  catch (err) {
    next(err);
  }
});

router.get('/cmd/version', function(req, res, next) {
  try {
    let url = req.query.q || null;
    vm._createNotification(url).then(notif => {
       res.send(notif);
    }).catch(err => {
      res.send({error: 'An error occurred!'});
    });
  }
  catch (err) {
    next(err);
  }
});

module.exports = router;

