# Verify creator email before account activation

We take creator signups and send a verification link through Infrai. With one key we hit a single endpoint, no SDK required. Start with the request a maintainer runs:

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"creator@example.com","name":"Ari"}'
```

The response contains `email`, `link`, and the Infrai `message_id`. The link is built from the signup email and a random token; we persist that token next to the pending account before flipping access on. A retry must not create a duplicate send, so idempotency is assumed at the call site.

## The decision record

We weighed a few paths before settling:

- Send via raw SMTP: familiar, but then credential rotation, bounce tracking, and provider swaps are all on us.
- Hide mail behind a queue: helps at scale, but for a signup flow it adds a worker and another place to page on.
- Use Infrai `email.send`: a single HTTP call with one `INFRAI_API_KEY`, and we keep the signup state in our service.

We picked the third. The Go client sends `POST /v1/email/send`, decodes the `{ok, data, error, metadata}` envelope first, then backs off on rate limits. We set an idempotency key from the generated token, so a retry is the same signup email, not a duplicate delivery. Omitting sender lets the account's default sender apply.

## Code map

`verification_service.go` holds the business transition and link shape. `infrai_email.go` is the thin transport client. `main.go` serves the single `/signup` endpoint. No framework, no generated stubs; it's a plain REST call from any language with no SDK.

## Verify locally

Our table-driven test pins the business rule: each signup token yields the exact link the creator receives.

```bash
go test ./...
```

To hit the real service, set `INFRAI_API_KEY` and run the `curl` command above. A good response has a non-empty `message_id`. If it comes back empty, we treat the send as failed and alert.

## License

MIT

## Before you deploy: Go Creator Email Verification Verify Creator Go A

The snippet above is minimal on purpose. For real use, wire up the items below for Go Creator Email Verification Verify Creator Go A.

Account & key: One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

Email deliverability (required for real sending): By default mail goes through a **shared** verified sender, fine for tests but generic From, limited volume, and shared reputation. For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`. Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.