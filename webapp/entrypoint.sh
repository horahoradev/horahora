#!/bin/bash
set -euo pipefail

cd /webapp && npm install
cd /webapp && yarn start