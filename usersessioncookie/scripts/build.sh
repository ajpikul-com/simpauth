#!/bin/bash
if [ "$1" == "test" ]; then
  npx eslint test/ && npx tsc -p tsconfig.test.json
else
  npx eslint src/ && npx tsc -p tsconfig.prod.json
fi
