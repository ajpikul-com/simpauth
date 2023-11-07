#!/bin/bash
npx eslint src/ && npx tsc -p tsconfig.prod.json
