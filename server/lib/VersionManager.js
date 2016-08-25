const request = require('request');
const Promise = require('bluebird');

const HOST = process.env.VERSION_HOST;
const VERSION = process.env.VERSION_ROUTE;

class VersionManager
{

  constructor(host, path)
  {
    this._host = host || HOST;
    this._version = path || VERSION;
  }

  _fetchVersion(url)
  {
    url = url || `https://${this._host}/${this._version}`;

    return new Promise(function(resolve, reject) {
      request(url, (err, resp, data) => {

        if (err) {
          return reject(err);
        }

        return resolve(data);
      });

    });
  }

  _getTimeLapsed(data)
  {
    const vTime = new Date(data.timestamp);
    const currentTime = new Date();

    // calc time passed
    const diffTime = currentTime - vTime;
    let resultString = '';

    if (diffTime > 60e3) {
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

  _createNotification(url)
  {
    return new Promise((resolve, reject) => {
      this._fetchVersion(url)
        .then(data => {
          const vData = JSON.parse(data);
          const sinceTimeString = this._getTimeLapsed(vData);
          const deployedDate = this._formatDate(new Date(vData.timestamp));
          resolve({
            response_type: 'in_channel',
            text: `Version: ${vData.version} deployed ${sinceTimeString}`,
            attachments: [
            {
              title: `Version ${vData.version}`,
              title_link: url || `https://${this._host}/${this._version}`,
              text: `New version deployed ${deployedDate}`
            }
            ]
          });
        }).catch(err => {
          reject(err);
        });
    });
  }

  sendNotification()
  {
    this._createNotification().then(payload => {
      //TODO send data to slack via https
    });
  }

}

module.exports = VersionManager;

