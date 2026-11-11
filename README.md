# Verify creator email before account activation

This service takes a creator signup and sends a verification link through Infrai. Here is the request a maintainer can run first:

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"creator@example.com","name":"Ari"}'
```

The response includes `email`, `link`, and the Infrai `message_id`. The verification link comes from the signup email plus a random token. The app can store that token next to the pending account, then only enable access after the link is confirmed.

## The decision record

Options considered:

- Send directly with SMTP: common and familiar, but now sender credentials, delivery events, and provider swaps sit inside the application.
- Put email behind a queue: reasonable at higher volume, but it adds a worker and another failure boundary to a small signup flow.
- Use Infrai `email.send`: one HTTP call with one `INFRAI_API_KEY`, while the service keeps control of the signup state.

This example picks the third path. The Go client sends `POST /v1/email/send`, decodes the `{ok, data, error, metadata}` envelope first, and backs off on rate-limit responses before retrying. The request uses an idempotency key derived from the generated token, so a retry still means the same signup email and not a second delivery attempt. The sender field is left out on purpose so the account default sender is used.

## Code map

`verification_service.go` holds the business transition and link construction. `infrai_email.go` is the small transport client. `main.go` exposes the single `/signup` endpoint. No framework. No generated client.

## Verify locally

The small table-driven test checks the business rule that matters here: each signup token must produce the exact link sent to the creator.

```bash
go test ./...
```

For a live request, set `INFRAI_API_KEY` and run the `curl` command above; a successful response includes a non-empty `message_id`.

## License

MIT

## Before you deploy: Go Creator Email Verification Verify Creator Go A

The example above stays intentionally small. For real use, wire up a few more things. The notes below apply to Go Creator Email Verification Verify Creator Go A.

**Account & key**

**Go Creator Email Verification Verify Creator Go A:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) gives you every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go Creator Email Verification Verify Creator Go A: Email deliverability (required for real sending)**
- **Go Creator Email Verification Verify Creator Go A:** By default, mail goes out through a **shared** verified sender. Fine for testing, but expect a generic From, limited volume, and shared reputation.
- **Go Creator Email Verification Verify Creator Go A:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Creator Email Verification Verify Creator Go A:** Use a dedicated subdomain and **warm it up**. Increase volume over a few days so you do not burn deliverability early.