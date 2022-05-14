#!/bin/bash
set -euo pipefail

cp node_modules/antd/dist/antd.dark.min.css public/antd.dark.min.css && \
    cp node_modules/antd/dist/antd.min.css public/antd.min.css

# TODO pay attention to linting lol
cd /webapp && yarn start