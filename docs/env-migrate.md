# MIGRATE FROM OLD ENVIRONMENT VARIABLES SETUP

1. Copy example `.env` file
    ```sh
    cp ./configs/.env.example .env
    ```

2. If you've changed anything in `secrets.env.template`, move the values to the `.env` file.

3. From `docker-compose.yml.envs` copy the value of `RSA_KEYPAIR` into `JWT_KEYPAIR` in `.env` file. The value shoudl be double quoted there.

4. Delete `docker-compose.yml.envs` and `secrets.env.template`.

5. In the new system the compose config is configured by the values in the `.env` file.
