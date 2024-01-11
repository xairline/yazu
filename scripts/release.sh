set -e

if [ -z "$1" ]
  then
    echo "No argument supplied, example: ./scripts/release.sh patch"
    exit 1
fi


version=$(npm version $@ --no-git-tag-version) # Adjust as needed
sed -i '' "s/\"version\": \".*\"/\"version\": \"$version\"/" wails.json
rm -rf build/bin
wails build --platform darwin/universal
npx --yes create-dmg build/bin/yazu.app build/bin --overwrite || true
rm -rf build/bin/yazu_*.dmg
mv "build/bin/yazu 1.0.0.dmg" "build/bin/yazu_${version}.dmg"

conventional-changelog -p angular -i CHANGELOG.md -s

git commit -am "Update changelog"
git tag ${version}
npx --yes gh-release \
  --assets build/bin/yazu_${version}.dmg \
  -t ${version} \
  --prerelease \
  -y -c main \
  -n ${version} \
  -o xairline \
  -r yet-another-zibo-updater
git push && git push --tags