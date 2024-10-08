#!/bin/bash
set -e
set -x

CURRENT_BRANCH=$(git branch --show-current)
VERSION=$(echo 0.0.$(( $(git tag|grep -E '^v0\.0\.'|sed -re 's/^v0\.0\.//'|sort -n|tail -n 1) + 1 )))
export BUILD_VERSION=$VERSION


#(
#	cd internal/frontend/statelessdb-frontend
#	git pull
#)
#git add internal/frontend/statelessdb-frontend
#(
#	cd internal/frontend/statelessdb-frontend
#	git tag -a "v$VERSION" -m "Version $VERSION"
#	git push origin "v$VERSION"
#	git push
#)

git tag -a "v$VERSION" -m "Version $VERSION"
git push origin "v$VERSION"
git push

git checkout main
git merge "$CURRENT_BRANCH"
git push

git checkout v0.0
git merge main
git push

git checkout "$CURRENT_BRANCH"

set +x
printf -- "READY!\n\nSee https://github.com/hyperifyio/statelessdb/releases\n"
