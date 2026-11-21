# ADR 0001: Store money as fixed-scale micros

- Status: Accepted
- Date: 2026-11-16
- Plan step: 1.5

## Context

`money.Money` held an amount in each currency's minor units: cents for CAD, kobo for NGN, whole
francs for XOF, which has no minor unit. Reading a stored amount therefore required knowing that
currency's number of decimal places, taken from `money.MinorUnits`.

This caused two problems:

1. If ISO 4217 changed the minor units of a currency we serve, every amount already stored in that
   currency would be read at the wrong scale, with no error.
2. Any ratio between two currencies had to be rescaled for the pair. Dividing raw minor units for
   CAD to XOF gives 4.3 instead of 430.

## Decision

Store every amount as an `int64` count of micros, millionths of one whole unit, for all
currencies: `Money{Micros: 12_500_000, Currency: "CAD"}` is CAD 12.50. Database columns are
`bigint` micros. This follows the `amount_micros` convention used by the Google Ads API.

`money.MinorUnits` stays. It is used to validate that an amount has no more precision than its
currency allows, and to format amounts for display. It no longer affects how stored data is read.

## Consequences

- Stored amounts are self-describing. A change to `money.MinorUnits` affects validation and
  display of new data only; no rescaling migration is needed.
- EffectiveRate is a plain division of two micros values, with no per-pair scaling.
- Arithmetic and comparison stay on a single `int64`.
- Precision is 1e-6 of a unit, finer than any ISO 4217 currency (at most 4 decimal places).
- Range is about plus or minus 9.2 trillion whole units, enough for any single transfer in the
  currencies we serve.
- Values are less readable in raw form (CAD 100.00 is 100000000), so logs and errors label them
  as micros.

## Alternatives considered

- **Keep minor units and store each row's scale** in a `minor_units` column. This protects stored
  data but keeps per-pair scaling in every cross-currency calculation.
- **`google.type.Money` (units + nanos).** Same benefit, but two fields with carry and sign rules
  make arithmetic harder, and importing the genproto type would put a protobuf type in the domain
  model. We do not exchange data with Google APIs, so matching that format gains nothing.
- **Decimal types** (for example `shopspring/decimal`, or Postgres `NUMERIC`). Exact, but adds a
  dependency and moves arithmetic off native integers.
