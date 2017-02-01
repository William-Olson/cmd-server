
const slackErr = {
  Message: 'Uh oh: Bad Request =('
};

module.exports = (fn, info) => {
  const route = `${info.method.toUpperCase()}: '${info.fullPath}'`;
  const name = `${info.routeClass}.${info.handler}`;

  console.log(`[harness] wrapping routes ({ ${name}, ${route} })`);

  return async (req, res, next) => {

    try {
      console.log(`[harness] ${route}, ${name}()`);
      const ret = await fn(req, res);

      if (ret) {
        res.send(ret);
      }
    }
    catch (err) {
      next(slackErr);
    }

  };
};
