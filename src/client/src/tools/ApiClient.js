import config from '../config';
const BASE = config.frankServerApi || 'http://localhost:8080';

function apiClient() {
  return {
    fetch: function(url, options) {
      if(!options) {
        options = {};
      }

      return fetch(BASE + url, options).then(filterError).then(parseJSON);
    },
    get: function(url, options) {
      return this.fetch(url, options);
    },
    delete: function(url, options) {
      if(!options) {
        options = {};
      }
      options.method = 'DELETE';
      options.headers = {
        'Content-Type': 'application/json'
      };
      return this.fetch(url, options);
    },
    patch: function(url, data, options) {
      if(!options) {
        options = {};
      }
      options.method = 'PATCH';
      options.headers = {
        'Content-Type': 'application/json'
      };
      options.body = JSON.stringify(data);
      return this.fetch(url, options);
    },
    put: function(url, data, options) {
      if(!options) {
        options = {};
      }
      options.method = 'PUT';
      options.headers = {
        'Content-Type': 'application/json'
      };
      options.body = JSON.stringify(data);
      return this.fetch(url, options);
    },
    post: function(url, data, options) {
      if(!options) {
        options = {};
      }
      options.method = 'POST';
      options.headers = {
        'Content-Type': 'application/json'
      };
      options.body = JSON.stringify(data);
      return this.fetch(url, options);
    },
    postFile: function(url, files) {
      let data = new FormData();
      files.forEach(file => {
        data.append('file', file)
      })
      let options = {method: 'post', body: data};

      return this.fetch(url, options);
    }
  }
}

function parseJSON(res) {
  return res.json();
}


function filterError(res) {
  if (res.status < 200 || res.status > 300) {
    const contentType = res.headers.get('Content-Type');

    if (!~contentType.indexOf('json')) {
      return Promise.reject(new Error(res));
    }

    return res.json().then(json => {
      return json;
    }).then(json => {
      return Promise.reject(json);
    });
  }

  return res;
}

export default apiClient();