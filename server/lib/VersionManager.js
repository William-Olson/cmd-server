const rp = require('request-promise');
const moment = require('moment');

const SERVER = process.env.VERSION_SERVER;
const HOST = process.env.VERSION_HOST;
const VERSION = process.env.VERSION_ROUTE;

class VersionManager
{

  constructor(server, host, path)
  {
    this._server  = server || SERVER;
    this._host    = host   || HOST;
    this._version = path   || VERSION;
  }

  async _fetchVersion(url)
  {
    url = url || `https://${this._server}.${this._host}/${this._version}`;
    return await rp(url);
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

  async createSlugNotif(slug, host, path)
  {
    slug = slug || this._server;
    host = host || this._host;
    path = path || this._version;

    const data = await this._fetchVersion(`https://${slug}.${host}/${path}`);
    const vData = JSON.parse(data);
    const sinceTimeString = this._getTimeLapsed(new Date(vData.timestamp));
    const deployedDate = this._formatDate(new Date(vData.timestamp));

    return {
      response_type: 'in_channel',
      text: `Version ${vData.version} for _*${slug}*_ was built ${sinceTimeString} ago`,
      attachments: [
        {
          title: `Version ${vData.version}`,
          title_link: `https://${slug}.${host}`,
          text: `Build Date: ${deployedDate}`
        }
      ]
    };
  }

  async createNotification(url)
  {
    const data = await this._fetchVersion(url);
    const vData = JSON.parse(data);
    const sinceTimeString = this._getTimeLapsed(new Date(vData.timestamp));
    const deployedDate = this._formatDate(new Date(vData.timestamp));

    return {
      response_type: 'in_channel',
      text: `-> ${url || `https://${this._server}.${this._host}`}\n` +
            `Version ${vData.version} was built ${sinceTimeString} ago`,
      attachments: [
        {
          title: `Version ${vData.version}`,
          title_link: url || `https://${this._server}.${this._host}`,
          text: `Build Date: ${deployedDate}`
        }
      ]
    };
  }


}

module.exports = VersionManager;

