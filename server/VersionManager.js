const rp = require('request-promise');
const moment = require('moment');

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

  async createSlugNotif(slug, host, path)
  {
    slug = slug || this._server;
    host = host || this._host;
    path = path || this._versionPath;

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

