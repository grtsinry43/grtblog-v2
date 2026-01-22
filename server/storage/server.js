const http = require("http");
const net = require("net");
const fs = require("fs");
const path = require("path");
const { URL } = require("url");

const STATIC_ROOT = path.join(__dirname, "html");
const LISTEN_PORT = 5555;
const API_UPSTREAM = "http://127.0.0.1:8080";

const MIME_TYPES = {
  ".html": "text/html; charset=utf-8",
  ".js": "application/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".txt": "text/plain; charset=utf-8",
  ".svg": "image/svg+xml",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".gif": "image/gif",
  ".ico": "image/x-icon",
  ".webp": "image/webp",
  ".woff2": "font/woff2",
  ".woff": "font/woff",
  ".ttf": "font/ttf",
  ".map": "application/json; charset=utf-8",
};

function setHeaderIfAbsent(res, name, value) {
  if (!res.getHeader(name)) {
    res.setHeader(name, value);
  }
}

function safeResolve(root, requestPath) {
  const decoded = decodeURIComponent(requestPath);
  const normalized = path.normalize(decoded).replace(/^(\.\.(\/|\\|$))+/, "");
  const resolved = path.join(root, normalized);
  if (!resolved.startsWith(root)) {
    return null;
  }
  return resolved;
}

function chooseCompressedVariant(filePath, acceptEncoding) {
  const canBr = acceptEncoding.includes("br");
  const canGzip = acceptEncoding.includes("gzip");

  if (canBr && fs.existsSync(`${filePath}.br`)) {
    return { path: `${filePath}.br`, encoding: "br" };
  }
  if (canGzip && fs.existsSync(`${filePath}.gz`)) {
    return { path: `${filePath}.gz`, encoding: "gzip" };
  }
  return { path: filePath, encoding: null };
}

function serveFile(res, filePath, acceptEncoding) {
  const { path: finalPath, encoding } = chooseCompressedVariant(
    filePath,
    acceptEncoding
  );
  const ext = path.extname(filePath);
  const mime = MIME_TYPES[ext] || "application/octet-stream";

  setHeaderIfAbsent(res, "Content-Type", mime);
  if (encoding) {
    res.setHeader("Content-Encoding", encoding);
  }

  const stream = fs.createReadStream(finalPath);
  stream.on("error", () => {
    res.statusCode = 404;
    res.end("Not Found");
  });
  stream.pipe(res);
}

function handleStatic(req, res) {
  if (!fs.existsSync(STATIC_ROOT)) {
    res.statusCode = 500;
    res.end(`Storage directory not found: ${STATIC_ROOT}. Did you run 'make preview-isr'?`);
    return;
  }

  const url = new URL(req.url, "http://localhost");
  let requestPath = url.pathname;
  const wantsRoot = requestPath === "/";

  if (requestPath === "/") {
    requestPath = "/index.html";
  }

  serveCleanPath(req, res, url, requestPath, wantsRoot);
}

function serveCleanPath(req, res, url, requestPath, wantsRoot) {
  const resolved = safeResolve(STATIC_ROOT, requestPath);
  if (!resolved) {
    res.statusCode = 400;
    res.end("Bad Request");
    return;
  }

  fs.stat(resolved, (err, stat) => {
    if (!err && stat.isDirectory()) {
      const indexPath = path.join(resolved, "index.html");
      fs.stat(indexPath, (indexErr) => {
        if (indexErr) {
          res.statusCode = 404;
          res.end("Not Found");
          return;
        }
        serveFile(res, indexPath, req.headers["accept-encoding"] || "");
      });
      return;
    }

    if (err) {
      const withHtml = `${resolved}.html`;
      fs.stat(withHtml, (htmlErr) => {
        if (!htmlErr) {
          serveFile(res, withHtml, req.headers["accept-encoding"] || "");
          return;
        }
        if (wantsRoot) {
          const fallback = path.join(STATIC_ROOT, "posts", "index.html");
          fs.stat(fallback, (fallbackErr) => {
            if (fallbackErr) {
              res.statusCode = 404;
              res.end("Not Found");
              return;
            }
            serveFile(res, fallback, req.headers["accept-encoding"] || "");
          });
          return;
        }
        res.statusCode = 404;
        res.end("Not Found");
      });
      return;
    }

    serveFile(res, resolved, req.headers["accept-encoding"] || "");
  });
}


function handleProxy(req, res) {
  const upstreamUrl = new URL(req.url, API_UPSTREAM);
  const upstreamOptions = {
    protocol: upstreamUrl.protocol,
    hostname: upstreamUrl.hostname,
    port: upstreamUrl.port,
    method: req.method,
    path: upstreamUrl.pathname + upstreamUrl.search,
    headers: {
      ...req.headers,
      host: upstreamUrl.host,
      "x-forwarded-for": req.socket.remoteAddress,
      "x-forwarded-proto": "http",
    },
  };

  const upstreamReq = http.request(upstreamOptions, (upstreamRes) => {
    res.writeHead(upstreamRes.statusCode || 502, upstreamRes.headers);
    upstreamRes.pipe(res);
  });

  upstreamReq.on("error", () => {
    res.statusCode = 502;
    res.end("Bad Gateway");
  });

  req.pipe(upstreamReq);
}

const server = http.createServer((req, res) => {
  if (req.url && req.url.startsWith("/api/v2")) {
    handleProxy(req, res);
    return;
  }
  handleStatic(req, res);
});

server.on("upgrade", (req, socket, head) => {
  if (!req.url || !req.url.startsWith("/api/v2")) {
    socket.destroy();
    return;
  }

  const upstream = net.connect(8080, "127.0.0.1", () => {
    const headers = [
      `GET ${req.url} HTTP/1.1`,
      `Host: 127.0.0.1:8080`,
      "Connection: Upgrade",
      "Upgrade: websocket",
    ];
    for (const [key, value] of Object.entries(req.headers)) {
      if (value === undefined) {
        continue;
      }
      if (Array.isArray(value)) {
        for (const entry of value) {
          headers.push(`${key}: ${entry}`);
        }
      } else {
        headers.push(`${key}: ${value}`);
      }
    }
    headers.push("\r\n");
    upstream.write(headers.join("\r\n"));
    if (head && head.length) {
      upstream.write(head);
    }
    socket.pipe(upstream);
    upstream.pipe(socket);
  });

  upstream.on("error", () => {
    socket.destroy();
  });
});

server.listen(LISTEN_PORT, () => {
  console.log(`Static server listening on http://127.0.0.1:${LISTEN_PORT}`);
});
