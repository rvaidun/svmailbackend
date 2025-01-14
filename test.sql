CREATE TABLE public.users (
	access_token varchar NULL,
	token_type varchar NULL,
	refresh_token varchar NULL,
	expiry integer NULL,
	email varchar NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (email)
);
