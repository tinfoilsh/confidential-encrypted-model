# Confidential Encrypted Model (Attested Key Release)

Example of securely running an encrypted model on Tinfoil and delivering keys to the enclave after successful attestation.

- **Decryption Service**: HTTP server that receives the encryption key to decrypt a model
- **Key Loader**: Client tool to verify an enclave's attestation and if successful, deliver the decryption key to the enclave

## Quick Start

### Encrypt a model

Run `encrypt.sh` or follow the encryption instructions here: https://github.com/redhat-et/coco-inferencing

### Attesting and providing the decryption key

```bash
go run ./keyloader --enclave encrypted-model.inf6.tinfoil.sh --repo tinfoilsh/confidential-encrypted-model --key private.pem
```
