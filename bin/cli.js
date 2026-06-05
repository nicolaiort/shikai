#!/usr/bin/env node
const { existsSync } = require('fs');
const { join } = require('path');
const { spawn } = require('child_process');

const ext = process.platform === 'win32' ? '.exe' : '';
const binary = join(__dirname, '..', 'dist', `shikai${ext}`);
const installScript = join(__dirname, '..', 'install.js');

function run() {
  const child = spawn(binary, process.argv.slice(2), { stdio: 'inherit' });
  child.on('exit', (code) => process.exit(code));
  child.on('error', () => {
    console.error(
      'shikai binary not found. Try running the install script manually:',
    );
    console.error(`  node ${installScript}`);
    process.exit(1);
  });
}

if (!existsSync(binary)) {
  const { fork } = require('child_process');
  const child = fork(installScript, [], { stdio: 'inherit' });
  child.on('exit', (code) => {
    if (code !== 0) {
      process.exit(code);
    }
    run();
  });
} else {
  run();
}
