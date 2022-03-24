BUILD_FILE_NAME = incognito-cli
UNIX_INSTALL_FILE = ./scripts/install_unix.sh
UNIX_UNINSTALL_FILE = ./scripts/uninstall_unix.sh

linux:
	bash $(UNIX_INSTALL_FILE) -n $(BUILD_FILE_NAME) -a

macos:
	bash $(UNIX_INSTALL_FILE) -n $(BUILD_FILE_NAME) -a

clean:
	bash $(UNIX_UNINSTALL_FILE) -n $(BUILD_FILE_NAME)