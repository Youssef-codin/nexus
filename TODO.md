## Stripe Top-Up Flow

### `POST /wallet/topup`

- [ ] Validate request body (amount in piastres, min value)
- [ ] Get wallet from DB using authenticated user ID
- [ ] Create a pending transaction in DB (type: `credit`, status: `pending`)
- [ ] Create Stripe Payment Intent (amount, currency: `egp`, metadata: `transaction_id`)
- [ ] Return `client_secret` and `transaction_id` to client

### `POST /webhook/stripe`

- [ ] Read raw request body
- [ ] Verify Stripe webhook signature using `stripe.ConstructEvent`
- [ ] Switch on event type:
  - `payment_intent.succeeded` → update transaction to `completed`, increment wallet balance
  - `payment_intent.payment_failed` → update transaction to `failed`
- [ ] Return 200 immediately (Stripe retries if it doesn't get 200)

### Queries Needed (sqlc)

- [ ] `CreateTransaction`
- [ ] `UpdateTransactionStatus`
- [ ] `GetTransactionByID`
- [ ] `IncrementWalletBalance`

### Misc

- [ ] Add `STRIPE_SECRET_KEY` and `STRIPE_WEBHOOK_SECRET` to `.env`
- [ ] Store `transaction_id` in Stripe metadata so you can look it up in the webhook
