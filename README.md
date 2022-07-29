# Gin-LSAT Proxy

An implementation of [gin-lsat](https://github.com/DhananjayPurohit/gin-lsat) middleware to demonstrate serving of static files and creating a paywall for paid resources.

## Steps to run:-
1. Clone the repo.
2. In the repo folder, run the command:-
    ```shell
    go get .
    ```
3. Create `.env` file (refer `.env_example`) and configure `LND_ADDRESS` and `MACAROON_HEX` for LND client or `LNURL_ADDRESS` for LNURL client, `LN_CLIENT_TYPE` (out of LND, LNURL) `ROOT_KEY` (for minting macaroons) and `PORT`. 
4. To start the server, run:-
    ```shell
    go run .
    ```
