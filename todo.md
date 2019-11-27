## Inserts??
- Split the code, it's not happening 

## Large loop-like queries
I need to perform a query to apply a marketing attribution rule. In my mind this means finding all conversions for a KPI, then looping over each customer journey to apply the rule. Is there any alternative?
- one action at a time
- transform to a new table
- ETL

## Multiple users
Ask Jordan - how can I allow users to inject their own data sources into Presto? Catalogs? Is that secure?

## Allow adding catalogs at runtime
We need to be able to add read presto DBs at runtime for tracks and billing