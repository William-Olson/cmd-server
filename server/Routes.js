const VersionManager = require('./VersionManager');
const ClientManager = require('./ClientManager');

module.exports = class Routes {

  constructor({ harness })
  {
    this._vm = new VersionManager();
    this._cm = new ClientManager();

    harness.get('/', this.getRoot);
    harness.get('/version', this.getVersion);
    harness.get('/client-slugs/:token', this.getClientManagerSlugs);
    harness.post('/slug-version', this.getSlugVersion);
  }

  /*
    Root route
  */
  async getRoot()
  {
    return { ok: true };
  }

  /*
    Get the client manager data for a token
  */
  async getClientManagerSlugs(req)
  {
    const token = req.params.token;
    await this._cm.readInClients();
    return (this._cm.getSlugs(token) || []);
  }

  /*
    General Version command route
  */
  async getVersion(req)
  {
    let url = req.query.q || null;
    return await this._vm.createNotification(url);
  }

  /*
    Get slug version
   */
  async getSlugVersion(req)
  {
    const cmd = req.body.command;
    let slug = req.body.text;
    const token = req.body.token;

    if (!cmd.includes('/version')) {
      throw new Error('Bad Command');
    }

    // handle multi param requests
    if (slug.includes(' ')) {
      const slugs = slug.trim().split(' ');

      if (slugs.length > 1) {
        return await this._vm.createMultiSlugNotif(token, slugs);
      }

      slug = slugs[0];
    }

    return await this._vm.createSlugNotif(token, slug);
  }
};
