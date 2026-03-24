all:
	go install ./cmd/area
	systemctl --user stop podman.socket
	systemctl --user disable podman.socket
	rm -f /run/user/$(id -u)/podman/podman.sock
	rmdir /run/user/$(id -u)/podman 2>/dev/null || true
