#!/usr/bin/env node
const { createWriteStream, existsSync, mkdirSync, chmodSync } = require('fs');
const { get } = require('https');
const { join } = require('path');
const { platform, arch } = require('os');

const { version } = require('./package.json');
const BINARY_DIR = join(__dirname, 'dist');

const platformMap = {
  linux: 'linux',
  darwin: 'darwin',
  win32: 'windows',
};

const archMap = {
  x64: 'amd64',
  arm64: 'arm64',
};

function fetch(url) {
  return new Promise((resolve, reject) => {
    const req = get(url, (res) => {
      if (
        res.statusCode >= 300 &&
        res.statusCode < 400 &&
        res.headers.location
      ) {
        resolve(fetch(res.headers.location));
        return;
      }
      if (res.statusCode !== 200) {
        reject(
          new Error(`Download failed with HTTP ${res.statusCode} from ${url}`),
        );
        return;
      }
      resolve(res);
    });
    req.on('error', reject);
    req.end();
  });
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest);
    fetch(url)
      .then((res) => {
        res.pipe(file);
        file.on('finish', () => {
          file.close();
          resolve();
        });
        file.on('error', (err) => {
          file.close();
          reject(err);
        });
      })
      .catch(reject);
  });
}

async function install() {
  const os = platformMap[platform()];
  const cpu = archMap[arch()];

  if (!os || !cpu) {
    console.error(
      `Unsupported platform: ${platform()} ${arch()}. ` +
        'shikai is available for linux, darwin, and windows on amd64 and arm64.',
    );
    process.exit(1);
  }

  const ext = os === 'windows' ? '.exe' : '';
  const binaryName = `shikai-${os}-${cpu}${ext}`;
  const url = `https://github.com/nicolaiort/shikai/releases/download/v${version}/${binaryName}`;
  const dest = join(BINARY_DIR, `shikai${ext}`);

  if (!existsSync(BINARY_DIR)) {
    mkdirSync(BINARY_DIR, { recursive: true });
  }

  if (existsSync(dest)) {
    console.log(`shikai v${version} already installed.`);
    return;
  }

  console.log(`Downloading shikai v${version} for ${os}-${cpu}...`);
  try {
    await download(url, dest);
    if (os !== 'windows') {
      chmodSync(dest, 0o755);
    }
    console.log(`shikai v${version} installed successfully.`);
  } catch (err) {
    console.error(
      `Failed to download shikai v${version}: ${err.message}`,
    );
    console.error(
      `Download URL: ${url}`,
    );
    console.error(
      'Make sure the release exists. You can also install via "go install github.com/nicolaiort/shikai@latest".',
    );
    process.exit(1);
  }
}

install();
