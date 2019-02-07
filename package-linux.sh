#!/bin/sh
# Define colours
LIGHTGREEN='\033[1;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

printf "\n${LIGHTGREEN}Compiling MiningHQ Miner for Linux${NC}\n\n"
printf "${YELLOW}Building server installer${NC}\n"
cd cli
make clean; make build_linux
cd ..
printf "${YELLOW}Building service installer${NC}\n"
cd install-service
make clean; make build_linux
cd ..
printf "${YELLOW}Building miner service${NC}\n"
cd miner-service
make clean; make build_linux
cd ..
printf "${YELLOW}Building uninstaller${NC}\n"
cd uninstaller
make clean; make build_linux
cd ..
printf "${YELLOW}Building GUI${NC}\n"
cd gui
make build_linux
cd ..
printf "${GREEN}Compile completed${NC}\n"
printf "\n${LIGHTGREEN}Packaging MiningHQ Miner for Linux${NC}\n\n"
if [ -d "packages/linux" ]; then
  rm -Rf packages/linux
fi
mkdir packages/linux
cp cli/bin/mininghq-server-installer packages/linux
printf "${YELLOW}Added server installer${NC}\n"
cp install-service/bin/install-service packages/linux
printf "${YELLOW}Added service installer${NC}\n"
cp miner-service/bin/miner-service packages/linux
printf "${YELLOW}Added miner service${NC}\n"
cp uninstaller/bin/uninstall-mininghq packages/linux
printf "${YELLOW}Added uninstaller${NC}\n"
cp gui/bin/linux-amd64/'MiningHQ Miner Manager' packages/linux
printf "${YELLOW}Added GUI${NC}\n"
printf "${GREEN}All parts added${NC}\n"
printf "\n${LIGHTGREEN}Create package${NC}\n\n"
cd packages/linux
tar -zcvf MiningHQ-Miner.tar.gz *
# find . -type f ! -name "*.tar.gz" -exec rm -rf {} \;
printf "\n${LIGHTGREEN}Removed temporary files${NC}\n\n"
cd ..
printf "${GREEN}Package created, available in packages/linux${NC}\n"
