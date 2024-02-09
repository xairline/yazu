set -e

if [ -z "$1" ]
  then
    echo "No argument supplied, example: ./scripts/release.sh patch"
    exit 1
fi

version=$(npm version $@ --no-git-tag-version)
sed -i '' "s/\"version\": \".*\"/\"version\": \"$version\"/" wails.json
# Create a Go file with the version
echo "package main\n\n// AppVersion is the current version of the app\nconst AppVersion = \"$version\"" > appversion.go
git commit -am "Bump version to $version"
rm -rf build/bin
wails build --platform windows/amd64,darwin/universal,linux/amd64
npx --yes create-dmg@6.1.0 build/bin/yazu.app build/bin --overwrite || true
mv "build/bin/yazu 1.0.0.dmg" "build/bin/yazu_${version}.dmg"
mv "build/bin/yazu-amd64.exe" "build/bin/yazu_${version}.exe"

conventional-changelog -p angular -i CHANGELOG.md -s || true

git commit -am "Update changelog"
git tag -f ${version}
npx gh-release \
  --assets build/bin/yazu_${version}.dmg,build/bin/yazu_${version}.exe \
  -t ${version} \
  --prerelease \
  -y -c main \
  -n 1.0.0 \
  -o xairline \
  -b CHANGELOG.md \
  -r yazu
git push && git push -f --tags