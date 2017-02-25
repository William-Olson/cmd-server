const fs = require('mz/fs');
const path = require('path');

module.exports = class ClientManager {

  constructor()
  {
    this._clients = {};
    this._file = path.join(__dirname, '/clients.json');
  }

  getSlugs(token)
  {
    return this._clients[token];
  }

  async update(token, slugs = [])
  {
    if (!this._clients[token]) {
      this._clients[token] = slugs;
    }
    else {
      this._clients[token] = [
        ...new Set([ ...this._clients[token], ...slugs ])
      ];
    }
    await this.saveToFile();
  }

  async readInClients()
  {
    const exists = await fs.exists(this._file);
    if (exists) {
      const data = await fs.readFile(this._file, { encoding: 'utf8' });
      this._clients = JSON.parse(data);
    }
  }

  async saveToFile()
  {
    await fs.writeFile(this._file, JSON.stringify(this._clients));
  }

};