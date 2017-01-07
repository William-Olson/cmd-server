const rp = require('request-promise');

const SERVER = process.env.VERSION_SERVER;
const HOST = process.env.VERSION_HOST;
const VERSION = process.env.VERSION_ROUTE;

class VersionManager
{

  constructor(server, host, path)
  {
    this._server = server || SERVER;
    this._host = host || HOST;
    this._version = path || VERSION;
  }

  async _fetchVersion(url)
  {
    url = url || `https://${this._server}.${this._host}/${this._version}`;
    return await rp(url);
  }

  _getTimeLapsed(data)
  {
    const vTime = new Date(data.timestamp);
    const currentTime = new Date();

    // calc time passed
    const diffTime = currentTime - vTime;
    let resultString = '';

    if ((diffTime / 60) > 60e3) {
      const hours = Math.floor((diffTime / 60e3) / 60);
      resultString = `${hours} hours ago`;
    }
    else if (diffTime > 60e3) {
      const minutes = Math.floor(diffTime / 60e3);
      resultString = `${minutes} minutes ago`;
    }
    else {
      const sec = Math.floor(diffTime / 1e3);
      resultString = `${sec} seconds ago`;
    }

    return resultString;
  }

  _formatDate(date)
  {
    const MONTHS = ['January', 'February', 'March', 'April', 'May', 'June',
          'July', 'August', 'September', 'October', 'November', 'December'];
    const day = date.getDate();
    const month = MONTHS[date.getMonth()];
    const year = date.getFullYear();
    return `${month} ${day}, ${year}`;
  }

  async createSlugNotif(slug, path)
  {
    slug = slug || `${this._server}`;
    path = path || `${this._host}/${this._version}`;
    const data = await this._fetchVersion(`https://${slug}.${path}`);
    const vData = JSON.parse(data);
    const sinceTimeString = this._getTimeLapsed(vData);
    const deployedDate = this._formatDate(new Date(vData.timestamp));
    return {
      response_type: 'in_channel',
      text: `Version ${vData.version} for _*${slug}*_ was deployed ${sinceTimeString}`,
      attachments: [
        {
          title: `Version ${vData.version}`,
          title_link: `https://${slug}.${path}`,
          text: `Deployed ${deployedDate}`
        }
      ]
    };
  }

  async createNotification(url)
  {
    const data = await this._fetchVersion(url);
    const vData = JSON.parse(data);
    const sinceTimeString = this._getTimeLapsed(vData);
    const deployedDate = this._formatDate(new Date(vData.timestamp));
    return {
      response_type: 'in_channel',
      text: `-> ${url || `https://${this._host}/${this._version}`}\n` +
            `Version ${vData.version} was deployed ${sinceTimeString}`,
      attachments: [
        {
          title: `Version ${vData.version}`,
          title_link: url || `https://${this._host}/${this._version}`,
          text: `Deployed ${deployedDate}`
        }
      ]
    };
  }


}

module.exports = VersionManager;

