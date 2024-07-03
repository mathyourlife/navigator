
## Setting up your local environment:

Create an empty sqlite database.

```
sqlite3 db/data/navigator.db ""
```

Install npm

```
npm install react react-dom react-modal
npm install @mui/material @emotion/react @emotion/styled
```


## Run locally

In one terminal start the frontend react server in develop mode that will auto update on saved file changes.

```bash
bash ./scripts/run-dev-frontend.sh
```

In another terminal start the backend API server that will process the data requests from the frontend.

```bash
bash ./scripts/run-dev-backend.sh
```
