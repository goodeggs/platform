.PHONY: rekey

SHELL=/bin/bash -o pipefail

rekey:
	@rm -f .rekey
	@git-crypt unlock
	@git-crypt export-key .rekey
	@printf '%s # GITCRYPT_PASS' "$$(travis encrypt GITCRYPT_PASS="$$(base64 .rekey | tr -d '\n')")"
	@rm -f .rekey
