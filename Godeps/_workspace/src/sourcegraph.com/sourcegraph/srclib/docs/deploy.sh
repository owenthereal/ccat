#!/usr/bin/env bash

set -e

# Create site dir if it does not exist
mkdir -p site

# Compile the css file
sass theme/styles.scss:theme/styles.css

# Build the other parts of the site
python2 buildsite.py

# Sync site with S3 bucket
aws s3 sync site/ s3://srclib.org/

echo <<EOF
You are not finished! If you're done deploying to the site, you need
to invalidate CloudFront's files with the following command:

    s3cmd --cf-invalidate sync site/ s3://srclib.org/

This costs $0.005 per file invalidated (after the first 1000 files in
a month), which isn't a ton, but it's best to only run it when you
need to.
EOF
