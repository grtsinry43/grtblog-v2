# Federation TODO

- [ ] Add Ed25519 signing/verification support once key format is finalized.
- [ ] Enforce per-instance rate limiting with Redis keys.
- [ ] Add SSRF protections for resolver/client fetches.
- [ ] Persist well-known metadata snapshots into federation_instance.
- [ ] Wire handlers/services to use signer/verifier and cache.
- [ ] Add background sync worker for timeline/RSS.
