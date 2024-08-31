# Synapse matrix account self service

## Run the service

The service runs on port 8080

Required environment variables:

- `KEYCLOAK_HOST`, ex.: 'auth.aalen.space'
- `KEYCLOAK_REALM`, ex.: 'master'
- `KEYCLOAK_CLIENT_ID`, ex.: 'matrix-self-service'
- `SYNAPSE_HOST`, ex.: 'matrix.aalen.space'
- `SYNAPSE_ACCESS_TOKEN`
> Request this via synapse rest api, since if you do it via the web client this
> token will be invalidated on logout
- `SYNAPSE_DOMAIN`, ex.: 'aalen.space'
 
