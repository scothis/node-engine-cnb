#!/usr/bin/env bash
set -euo pipefail

TARGET_OS=${1:-linux}

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

echo "Target OS is $TARGET_OS"
echo -n "Creating buildpack directory..."
bp_dir="${PWD##*/}"_$(openssl rand -hex 12)
mkdir $bp_dir
echo "done"

echo -n "Copying buildpack.toml..."
cp buildpack.toml $bp_dir/buildpack.toml
echo "done"

if [ "${BP_REWRITE_HOST:-}" != "" ]; then
    sed -i -e "s|^uri = \"https:\/\/buildpacks\.cloudfoundry\.org\(.*\)\"$|uri = \"http://$BP_REWRITE_HOST\1\"|g" "$bp_dir/buildpack.toml"
fi

for b in $(ls cmd); do
    echo -n "Building $b..."
    GOOS=$TARGET_OS go build -o $bp_dir/bin/$b ./cmd/$b
    echo "done"
done
echo "Buildpack packaged into: $bp_dir"

pushd $bp_dir
    tar czvf ../nodejs-cnb.tgz *
popd

shasum="$(shasum -a 256 nodejs-cnb.tgz | cut -d ' ' -f1)"

rm -rf $bp_dir
rm -rf "/Users/pivotal/.buildpack-packager/cache"


script=$(cat <<EOF
require 'YAML'
file = '/Users/pivotal/workspace/nodejs-buildpack/manifest.yml'
m = YAML.load_file(file)
m['dependencies'][0]['sha256'] = '$shasum'
File.open(file, 'w') {|f| f.write m.to_yaml }
EOF
)
ruby -e "$script"