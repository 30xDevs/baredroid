#!/bin/bash
# _TOOLS=/opt/android-sdk-update-manager/build-tools/29.0.3
_AAPT=aapt2
adb shell pm list packages --user 0 | sed -e 's|^package:||' | sort >./packages_list.txt
_PMLIST=packages_list.txt
rm ./packages_list_with_names.json
_TEMP=$(echo $(adb shell mktemp -d -p /data/local/tmp) | sed 's/\r//')
mkdir -p packages
[ -f ${_PMLIST} ] || eval $(echo $(basename ${_PMLIST}) | tr '_' ' ') > ${_PMLIST}

# Initialize json
echo "[" > ./packages_list_with_names.json
while read -u 9 _line; do
    _package=${_line##*:}
    _apkpath=$(adb shell pm path ${_package} | sed -e 's|^package:||' | head -n 1)
    _apkfilename=$(basename "${_apkpath}")
    adb shell cp -f ${_apkpath} ${_TEMP}/copy.apk
    adb pull ${_TEMP}/copy.apk ./packages
    _name=$(${_AAPT} dump badging ./packages/copy.apk | sed -n 's|^application-label:\(.\)\(.*\)\1$|\2|p' )
#'
    echo "\t { \n \"name\": \"${_name}\",\n \"package\": \"${_package}\"\n }," >>./packages_list_with_names.json
done 9< ${_PMLIST}

# Remove last comma
sed -i '$ s/,$//' ./packages_list_with_names.json

# Finish json list
echo "]" >> ./packages_list_with_names.json

adb shell rm -rf $TEMP
