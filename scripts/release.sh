set -e
wails build --platform darwin/universal
npx --yes create-dmg build/bin/yazu.app build/bin --overwrite || true

version=$(npm version $@) # Adjust as needed
git tag version

rm -rf build/bin/yazu_*.dmg
mv build/bin/yazu\ 1.0.0.dmg build/bin/yazu_${version}.dmg

conventional-changelog -p angular -i CHANGELOG.md -s

git commit -am "Update changelog"

npx --yes gh-release \
  --assets build/bin/yazu_${version}.dmg \
  -t ${version} \
  --prerelease \
  -y -c main \
  -n ${version} \
  -o xairline \
  -r yet-anther-zibo-updater
git push && git push --tags