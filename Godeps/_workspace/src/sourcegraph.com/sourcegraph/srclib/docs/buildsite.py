#!/usr/bin/env python

import jinja2
import os
import re
import shlex
import sys

import mkdocs.build
from mkdocs.build import build
from mkdocs.config import load_config
from urllib2 import urlopen
import subprocess

def line_containing(lines, text):
  for i in range(len(lines)):
    if text.lower() in lines[i].lower():
      return i
  raise Exception("could not find {}".format(text))

# Wrap some functions to allow custom commands in markdown
convert_markdown_original = mkdocs.build.convert_markdown
def convert_markdown_new(source, **kwargs):

  def expand(match):
    args = shlex.split(match.groups()[0])

    # Import external markdown
    if args[0] == ".import":
      code = ""
      try: #Try as a URL
        code = urlopen(args[1]).read()
      except ValueError:  # invalid URL, try as a file
        code = open(args[1]).read()

      return code

    # Run a shell command
    elif args[0] == ".run":
      result = ""
      command = "$ " + match.groups()[0].replace(".run", "").strip()
      try:
        result = subprocess.check_output(args[1:],  stderr=subprocess.STDOUT)
      except subprocess.CalledProcessError, e:
        result = e.output
      return "```\n" + command + "\n" + result.strip() + "\n```"

    # Source code embeds
    elif args[0] == ".code" or args[0] == ".doc":
      code = ""
      try: #Try as a URL
        code = urlopen(args[1]).read()
      except ValueError:  # invalid URL, try as a file
        code = open("../" + args[1]).read()

      lines = code.splitlines()

      # Short hand for specifying a region
      if len(args) == 3:
        region = args[2]
        args[2] = "START " + region
        args.append("END " + region)

      if len(args) == 4:
        start = 1
        end = len(lines) - 1

        try:
          if args[2].isdigit(): start = int(args[2])
          else:
            start = line_containing(lines, args[2]) + 1

          if args[3].isdigit(): end = int(args[3])
          else: end = line_containing(lines, args[3]) + 1
        except Exception, e: # If line_containing fails
          print "Error: {}".format(e)
          print "  in {}".format(args[1])
          sys.exit(1)

        #TODO: Also allow regex matching

        lines = lines[start - 1:end]

      # Trim "OMIT" lines. Ignore "*/".
      lines = filter(lambda x: not x.strip().rstrip("*/").rstrip().lower().endswith("omit"), lines)

      # TODO: Trim leading and trailing empty lines

      if args[0] == ".code":
        lines.insert(0, "```go")
        lines.append("```")
      # else: # args[0] == ".doc"
      #   lines.insert(0, "\n")
      #   lines.insert("\n")
      return "\n".join(lines)

    # No matching logic
    else:
      return match.group(0)
  # Process an aritrary number of expansions.
  oldSource = ""
  while source != oldSource:
    oldSource = source
    source = re.sub("\[\[(.*)\]\]", expand, oldSource)

  return convert_markdown_original(source)

# Hotpatch in the markdown conversion wrapper
mkdocs.build.convert_markdown = convert_markdown_new




if __name__ == "__main__":
  # Build documentation
  config = load_config(options=None)
  build(config)

  # Load templates
  template_env = jinja2.Environment(loader = jinja2.FileSystemLoader(os.path.join(os.path.dirname(__file__), 'theme')))
  index_template = template_env.get_template('home.html')
  community_template = template_env.get_template('community.html')

  # Home page
  with open('site/index.html', 'w') as f:
    f.write(index_template.render(
      page="home"
    ))

  # Community page
  with open('site/community.html', 'w') as f:
    f.write(community_template.render(
      page="community"
    ))
