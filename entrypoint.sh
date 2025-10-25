#!/bin/sh
set -e

# If any of the arguments contains --url= or --url we assume load test mode
for arg in "$@"; do
    case "$arg" in
        --url|--url=*)
            exec ./loadtest "$@"
            ;;
    esac
done

exec ./server "$@"


