"use strict"

module.exports = (event, context) => {
    const { method, path } = event;

    if (method !== 'GET') {
      context.status(400).fail('Bad Request');
      return;
    }

    let headers = {
      'Content-Type': '',
    };
    if (/.*\.js/.test(path)) {
      headers['Content-Type'] = 'application/javascript';
    } else if (/.*\.css/.test(path)) {
      headers['Content-Type'] = 'text/css';
    } else if (/.*\.ico/.test(path)) {
      headers['Content-Type'] = 'image/x-icon';
    } else if (/.*\.json/.test(path)) {
      headers['Content-Type'] = 'application/json';
    } else if (/.*\.map/.test(path)) {
      headers['Content-Type'] = 'application/octet-stream';
    }
  
    let contentPath = `${__dirname}${path}`;
  
    if (!headers['Content-Type']) {
      contentPath = `${__dirname}/client/dist/index.html`;
    }
  
    fs.readFile(contentPath, (err, data) => {
      if (err) {
        context
          .headers(headers)
          .status(500)
          .fail(err);
  
        return;
      }
  
      let content = data.toString();

      context
        .headers(headers)
        .status(200)
        .succeed(content);
    });
}
