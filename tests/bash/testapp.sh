#!/bin/bash

# vim: filetype=sh:tabstop=2:shiftwidth=2:expandtab

readonly PROGNAME=$(basename $0)
readonly PROGDIR="$( cd "$(dirname "$0")" ; pwd -P )"
readonly DEFAULT_PROGDIR=/home/pinter/projects/github
readonly DEFAULT_GH_USERNAME=pinterb
readonly CURL_CMD=`which curl`
readonly JQ_CMD=`which jq`
readonly GIT_CMD=`which git`

readonly GH_API_BASE_URI=https://api.github.com

declare -r TRUE=0
declare -r FALSE=1

# Get to where we need to be.
cd $PROGDIR

# Globals overridden as command line arguments
PROJECT_DIRECTORY=$DEFAULT_PROGDIR
GH_USER=$DEFAULT_GH_USERNAME

usage()
{
  echo -e "\033[33mHere's how to show which repos you have forked on GitHub:"
  echo ""
  echo -e "\033[33m./$PROGNAME"
  echo -e "\t\033[33m-h --help"
  echo -e "\t\033[33m--user=$GH_USER (i.e. GitHub username)"
  echo -e "\t\033[33m--dir=$PROJECT_DIRECTORY (i.e. directory the forked repository will be cloned into)"
  echo -e "\033[0m"
}


parse_args()
{
  while [ "$1" != "" ]; do
    PARAM=`echo $1 | awk -F= '{print $1}'`
    VALUE=`echo $1 | awk -F= '{print $2}'`
    case $PARAM in
      -h | --help)
        usage
        exit
        ;;
      --user)
        GH_USER=$VALUE
        ;;
      --dir)
        PROJECT_DIRECTORY=$VALUE
        ;;
      *)
        echo -e "\033[31mERROR: unknown parameter \"$PARAM\""
        echo -e "\e[0m"
        usage
        exit 1
        ;;
    esac
    shift
  done

}


valid_args()
{

  if [ ! -d "$PROJECT_DIRECTORY" ]; then
    echo -e "\033[31mERROR: directory \"$PROJECT_DIRECTORY\" does not exist"
    echo -e "\e[0m"
    usage
    exit 1
  fi

}


show_forks()
{
  cd $PROJECT_DIRECTORY
  shopt -s dotglob
  find * -prune -type d | while read d; do
    if [ -d "$PROJECT_DIRECTORY/$d/.git" ]; then
      git -C $PROJECT_DIRECTORY/$d remote | grep -vq origin
      RETVAL=$?
      if [ $RETVAL -eq 0 ]; then
        my_remote=`git -C $PROJECT_DIRECTORY/$d remote | grep -v origin` 
        my_remote_url=`git -C $PROJECT_DIRECTORY/$d config --get remote.${my_remote}.url`
        echo "${d}: ${my_remote_url}"
      fi
    fi
  done

}


main()
{
  # Perform sanity check on command line arguments
  valid_args

  # display git repositories that have a remote other than "origin"
  show_forks
}


parse_args "$@"
main
exit 0
