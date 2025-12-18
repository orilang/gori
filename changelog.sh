#!/bin/bash

#
# Generate a Markdown-formatted changelog from merge commits.
#

set -o errexit

usage(){
  echo "$1 argument is mantatory"
  echo "Example: $0 init 0.0.1"
  echo "Example: $0 0.0.1 0.1.0"
}

if [ -n "$1" ]
then
  export PREVIOUS_TAG="$1"
else
  usage PREVIOUS_TAG 
  exit 1
fi
shift

if [ -n "$1" ]
then
  export NEW_TAG="$1"
else
  usage NEW_TAG 
  exit 1
fi
shift

MD_HEADER_SOURCE_NAME="changelog/.header/header.tmpl.md"
MD_HEADER_TARGET_NAME="changelog/.header/v${NEW_TAG}_header.md"
MD_BODY_NAME="changelog/.header/v${NEW_TAG}_body.md"
MD_NAME="changelog/v${NEW_TAG}.md"
TODAY_DATE=$(date +%Y-%m-%d)

#
# Regex to match the commit message. Creates capture groups for git
# author, and commit subject.
#
parse_git_log(){
  PARSE_RE='([^\|]+)\|(.*)'
  git log --first-parent --format='%aN|%s' ${1} |
  {
    output=""
    while read l
    do
      output="`echo main`_output"
      if [[ ${l} =~ ${PARSE_RE} ]]
      then
        author="${BASH_REMATCH[1]}"
        summary="${BASH_REMATCH[2]}"
        declare $output+=" * ${summary} (by @${author})\n"
      fi
    done

    echo -e "$main_output" > "${MD_BODY_NAME}"
  }
}

manage_header(){
  OS=$(uname -s | tr -s "A-Z" "a-z")
  if [[ ! -f "${MD_HEADER_TARGET_NAME}" ]]
  then
    cp "${MD_HEADER_SOURCE_NAME}" "${MD_HEADER_TARGET_NAME}"
    if [[ "${OS}" == "darwin" ]]
    then
      sed -i '' "s#NEW_TAG#${NEW_TAG}#g" "${MD_HEADER_TARGET_NAME}"
      sed -i '' "s#YYY-MM-DD#${TODAY_DATE}#g" "${MD_HEADER_TARGET_NAME}"
    else
      sed -i "s#NEW_TAG#${NEW_TAG}#g" "${MD_HEADER_TARGET_NAME}"
      sed -i "s#YYYY-MM-DD#${TODAY_DATE}#g" "${MD_HEADER_TARGET_NAME}"
    fi
  else  
    if [[ "${OS}" == "darwin" ]]
    then
      sed -E "s#([0-9]{4}-[0-9]{2}-[0-9]{2})#${TODAY_DATE}#g" "${MD_HEADER_TARGET_NAME}"
    else
      sed -i -E "s#([0-9]{4}-[0-9]{2}-[0-9]{2})#${TODAY_DATE}#g" "${MD_HEADER_TARGET_NAME}"
    fi
    echo ""
  fi
}

create_release(){
  cat "${MD_HEADER_TARGET_NAME}" > "${MD_NAME}"
  cat "${MD_BODY_NAME}" >> "${MD_NAME}"
  rm -f "${MD_BODY_NAME}"
}

if [[ "${PREVIOUS_TAG}" == "init" ]]
then
   parse_git_log ""
else
  parse_git_log "v${PREVIOUS_TAG}.."
fi

manage_header
create_release
cat "${MD_NAME}"
