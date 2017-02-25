const rp = require('request-promise');
const moment = require('moment');
const ClientManager = require('./ClientManager');

const SERVER = process.env.VERSION_SERVER;
const HOST = process.env.VERSION_HOST;
const VERSION = process.env.VERSION_ROUTE;

class VersionManager
{

  constructor(server, host, path)
  {
    this._server      = server || SERVER;
    this._host        = host   || HOST;
    this._versionPath = path   || VERSION;

    this._clientManager = new ClientManager();
  }

  async _fetchVersion(url)
  {
    url = url || `https://${this._server}.${this._host}/${this._versionPath}`;
    return await rp({ url, json: true });
  }

  _getTimeLapsed(date)
  {
    const vTime = moment(date);
    const currentTime = moment(new Date());
    return moment.duration(currentTime.diff(vTime)).humanize();
  }

  _formatDate(date)
  {
    return moment(date).format('MMMM Do, YYYY');
  }

  async _getVersionInfo(url)
  {
    const payload = await this._fetchVersion(url);
    return {
      version: payload.version,
      sinceTime: this._getTimeLapsed(new Date(payload.timestamp)),
      buildDate: this._formatDate(new Date(payload.timestamp))
    };
  }

  async createMultiSlugNotif(token, slugs, host, path)
  {
    if (!slugs || !slugs.length || !token) {
      throw new Error('token and slugs required (slugs must be an array)');
    }

    await this._clientManager.update(token, slugs);

    host = host || this._host;
    path = path || this._versionPath;
    const attachments = [];

    for (let slug of slugs) {
      const url = `https://${slug}.${host}/${path}`;
      const { version, sinceTime, buildDate } = await this._getVersionInfo(url);
      const capSlug = `${slug[0].toUpperCase()}${slug.substr(1, slug.length - 1)}`;
      attachments.push({
        title: `${capSlug} is running ${version}`,
        title_link: `https://${slug}.${host}`,
        text: `Image built ${buildDate} (${sinceTime} ago)`
      });
    }

    return {
      response_type: 'in_channel',
      text: 'Listing versions',
      attachments
    };
  }

  async createSlugNotif(token, slug, host, path)
  {
    slug = slug || this._server;
    host = host || this._host;
    path = path || this._versionPath;

    if (!token) {
      throw new Error('Missing param token (required)');
    }

    if (slug === 'all') {
      const slugs = this._clientManager.getSlugs(token);
      return await this.createMultiSlugNotif(token, slugs, host, path);
    }

    await this._clientManager.update(token, [ slug ]);

    const url = `https://${slug}.${host}/${path}`;
    const { version, sinceTime, buildDate } = await this._getVersionInfo(url);
    const capitalizedSlug = `${slug[0].toUpperCase()}${slug.substr(1, slug.length - 1)}`;

    return {
      response_type: 'in_channel',
      text: `_*${capitalizedSlug}*_ is running version *${version}*`,
      attachments: [
        {
          title: `${slug}.${host}`,
          title_link: `https://${slug}.${host}`,
          text: `Image built ${buildDate} (${sinceTime} ago)`
        }
      ]
    };
  }

  async createNotification(url)
  {
    const { version, sinceTime, buildDate } = await this._getVersionInfo(url);

    return {
      response_type: 'in_channel',
      text: `${url || `https://${this._server}.${this._host}`}\n` +
            `Version ${version} was built ${sinceTime} ago`,
      attachments: [
        {
          title: `Version ${version}`,
          title_link: url || `https://${this._server}.${this._host}`,
          text: `Build Date: ${buildDate}`
        }
      ]
    };
  }


}

module.exports = VersionManager;

