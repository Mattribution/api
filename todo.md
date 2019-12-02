## Inserts??
- Split the code, it's not happening 

## Large loop-like queries
I need to perform a query to apply a marketing attribution rule. In my mind this means finding all conversions for a KPI, then looping over each customer journey to apply the rule. Is there any alternative?
- one action at a time
- transform to a new table
- ETL

## Multiple users
Ask Online - how can I allow users to inject their own data sources into Presto? Catalogs? Is that secure?

## Allow adding catalogs at runtime
We need to be able to add read presto DBs at runtime for tracks and billing

## Queries
- MostActiveCampaigns  
Finds the campaign names with most activity in general  
`select campaign_name, count(*) count from tracks group by 1 order by 2 desc;`
- MostActiveMediums
Finds the mediums with most activity
`select campaign_name, count(*) count from tracks group by 1 order by 2 desc;`
- FirstTouch  
Attributes conversions via first touch. This needs to find conversions, then get the first relevant touch point for the journey.