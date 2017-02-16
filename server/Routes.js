const VersionManager = require('./VersionManager');

module.exports = class CmdRoutes {

  constructor({ harness })
  {
    this._vm = new VersionManager();

    harness.get('/', this.getRoot);
    harness.get('/version', this.getVersion);
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
    General Version command route
  */
  async getVersion(req)
  {
    let url = req.query.q || null;
    return await this._vm.createNotification(url);
  }

  /**
   * get slug version
   */
  async getSlugVersion(req)
  {
    const cmd = req.body.command;
    let slug = req.body.text;

    if (!cmd.includes('/version')) {
      throw new Error('Bad Command');
    }

    // handle multi param requests
    if (slug.includes(' ')) {
      const slugs = slug.trim().split(' ');

      if (slugs.length > 1) {
        return await this._vm.createMultiSlugNotif(slugs);
      }

      slug = slugs[0];
    }

    return await this._vm.createSlugNotif(slug);
  }
};
