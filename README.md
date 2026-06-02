# Artist Insights

![Screenshot](/showcase.png)

This is a full-stack project to visualize and provide reports and insights on the Last.fm song listening history dataset and to provide insights for artists so that they can infer and improve on their fan base.

## Overview

The repository is built mainly by importing the Last.fm dataset and then adding it to ClickHouse, which stores the listening data points for each user (which song they listen to and when). ClickHouse is connected to it as a core server, which exposes its reports as an API and a React frontend using Recharts to present this to the user.

## Prerequisites

Get the dataset from http://ocelma.net/MusicRecommendationDataset/lastfm-1K.html and add it to the dataset folder.

## Getting Started

Provision ClickHouse and the Go Server using the makefile commands.

```
make provision
```

Finally, start the frontend server using. 

```
npm run dev
```

And head to http://localhost:5173/ to view the Insights page where you can search for each artist and view their reports.


## Future Scope

- Extend to simulate real-time insights like current plays and users.
- More artist-specific reports.
- Observability via oTel for both app and ClickHouse metrics.

## References

- http://ocelma.net/MusicRecommendationDataset/lastfm-1K.html
- https://clickhouse.com/docs/install/docker
- https://clickhouse.com/docs/engines/table-engines/mergetree-family
