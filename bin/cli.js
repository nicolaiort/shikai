#!/usr/bin/env node
const { join } = require('path');
const { spawn } = require('child_process');

const ext = process.platform === 'win32' ? '.exe' : '';
const binary = join(__dirname, '..', 'dist', `shikai${ext}`);

const child = spawn(binary, process.argv.slice(2), { stdio: 'inherit' });
child.on('exit', (code) => process.exit(code));
