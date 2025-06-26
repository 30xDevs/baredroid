# Baredroid
ADB abstraction program to configure Android devices.

## Goals
1. Install and Uninstall defined packages from .baredroid file.
2. Change settings according to .baredroid file.

## Potential application commands
- bdroid tools install {TOOL_NAME | None}
  - Will pull aapt (maybe more tools) and add to PATH

## Common commands
List packages
* adb shell pm list packages -s (-3 for user-installed packages, include `| cut -f 2 -d ":"` to remove the `package:` prefix)
  
Remove packages
* adb uninstall com.example.app