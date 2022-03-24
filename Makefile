BUILD_FILE_NAME = incognito-cli
UNIX_INSTALL_FILE = ./scripts/install_unix.sh
UNIX_UNINSTALL_FILE = ./scripts/uninstall_unix.sh

install:
	chmod +x $(UNIX_INSTALL_FILE) && bash $(UNIX_INSTALL_FILE) -n $(BUILD_FILE_NAME) -a

clean:
	chmod +x $(UNIX_UNINSTALL_FILE) && bash $(UNIX_UNINSTALL_FILE) -n $(BUILD_FILE_NAME)