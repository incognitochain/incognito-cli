BUILD_FILE_NAME = incognito-cli
LINUX_INSTALL_FILE = ./scripts/install_linux.sh
LINUX_UNINSTALL_FILE = ./scripts/uninstall_linux.sh
MACOS_INSTALL_FILE = ./scripts/install_macos.sh
MACOS_UNINSTALL_FILE = ./scripts/uninstall_macos.sh

linux:
	chmod +x $(LINUX_INSTALL_FILE) && bash $(LINUX_INSTALL_FILE) -n $(BUILD_FILE_NAME) -a

macos:
	chmod +x $(MACOS_INSTALL_FILE) && bash $(MACOS_INSTALL_FILE) -n $(BUILD_FILE_NAME) -a

#clean:
#	chmod +x $(UNIX_UNINSTALL_FILE) && bash $(UNIX_UNINSTALL_FILE) -n $(BUILD_FILE_NAME)