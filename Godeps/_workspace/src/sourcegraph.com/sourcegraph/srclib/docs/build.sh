#!/usr/bin/env bash

# Create site dir if it does not exist
mkdir -p site

# Compile the css file
sass theme/styles.scss:theme/styles.css

# Build the other parts of the site
python buildsite.py
