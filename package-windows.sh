#!/bin/sh
# Define colours
LIGHTGREEN='\033[1;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

printf "\n${LIGHTGREEN}Compiling MiningHQ Miner for Windows${NC}\n\n"
# printf "${YELLOW}Building service installer${NC}\n"
# cd install-service
# make clean; make build_windows
# cd ..
printf "${YELLOW}Building miner service${NC}\n"
cd miner-service
make clean; make build_windows
cd ..
printf "${YELLOW}Building uninstaller${NC}\n"
cd uninstaller
make clean; make build_windows
cd ..
printf "${YELLOW}Building GUI${NC}\n"
cd gui
make build_windows
cd ..
printf "${GREEN}Compile completed${NC}\n"
printf "\n${LIGHTGREEN}Packaging MiningHQ Miner for Windows${NC}\n\n"
if [ ! -d "packages" ]; then
  mkdir packages
fi
if [ -d "packages/windows" ]; then
  rm -Rf packages/windows
fi
mkdir packages/windows
mkdir packages/windows/tools
# cp install-service/bin/install-service.exe packages/windows/tools
# printf "${YELLOW}Added service installer${NC}\n"
cp miner-service/windows/run-as-service.bat packages/windows/tools
printf "${YELLOW}Added run-as-service${NC}\n"
cp miner-service/bin/miner-service.exe packages/windows/tools
printf "${YELLOW}Added miner service${NC}\n"
cp uninstaller/bin/uninstall-mininghq.exe packages/windows/tools
printf "${YELLOW}Added uninstaller${NC}\n"
cp gui/bin/windows-amd64/'MiningHQ Miner Manager.exe' packages/windows/'MiningHQ Miner Installer.exe'
printf "${YELLOW}Added GUI${NC}\n"
printf "${GREEN}All parts added${NC}\n"
printf "\n${LIGHTGREEN}Create package${NC}\n\n"
# cd packages/windows
# zip MiningHQ-Miner.zip *
# find . -type f ! -name "*.zip" -exec rm -rf {} \;
# printf "\n${LIGHTGREEN}Removed temporary files${NC}\n\n"
# cd ..
printf "${GREEN}Package created, available in packages/windows${NC}\n"
