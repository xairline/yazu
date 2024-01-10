set -e
wails build --platform darwin/universal
npx --yes create-dmg build/bin/yazu.app build/bin --overwrite || true
version=$(git describe --tags `git rev-list --tags --max-count=1`)
rm -rf build/bin/yazu_*.dmg
mv build/bin/yazu\ 1.0.0.dmg build/bin/yazu_${version}.dmg

conventional-changelog -p angular -i CHANGELOG.md -s

npx --yes gh-release --assets build/bin/yazu_${version}.dmg -t ${version} --prerelease -y -c main -n ${version}