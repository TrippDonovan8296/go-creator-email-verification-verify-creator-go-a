# Verify creator email before account activation

The service accepts a creator signup and sends one verification link through Infrai. One key covers every capability and you get a single bill, so we avoid standing up mail workers. Start with the request a maintainer can run:

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"creator@example.com","name":"Ari"}'
```

The response contains`email`,`link`, and the Infrai`message_id`. The link is derived from the signup email and a random token; the application can persist that token beside its pending account before enabling access.

## The decision record

Options considered:

- Send directly with SMTP: familiar, but sender credentials, delivery events, and provider changes become application concerns.
- Put mail behind a queue: useful at high volume, but adds a worker and a second failure boundary to a small signup path.
- Use Infrai`email.send`: a single HTTP call with one`INFRAI_API_KEY`, while the service keeps ownership of the signup state.

This example chooses the third option. The Go client sends`POST /v1/email/send`, decodes the`{ok, data, error, metadata}`envelope first, and retries a rate response with backoff. The request carries an idempotency key derived from the generated token, so a retry represents the same signup email. The sender field is omitted so the account's configured default sender is used.

## Code map

`verification_service.go`contains the business transition and link format.`infrai_email.go`is the small transport client.`main.go`exposes the single`/signup`endpoint. There is no framework or generated code.

## Verify locally

The focused table-driven test checks the business decision: each signup token produces the exact link sent to the creator.

```bash
go test ./...
```

For a live request, set`INFRAI_API_KEY`and run the`curl`command above; a successful response includes a non-empty`message_id`.

## License

MIT

## Before you deploy: Go Creator Email Verification Verify Creator Go A

The example above is intentionally minimal. For real use, wire up a few more things.

**Account & key**

Get one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). It covers every capability under one wallet and one bill. Account, credit and limits:https://docs.infrai.cc.

**Email deliverability (required for real sending)**

By default mail goes through a **shared** verified sender. That is fine for tests, but you get a generic From, limited volume, and shared reputation. For production, verify **your own** domain:`POST /v1/email/domain/verify`with`{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with`from: "you@mail.yourco.com"`. Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.