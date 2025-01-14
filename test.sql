CREATE TABLE public.users (
	access_token varchar NULL,
	token_type varchar NULL,
	refresh_token varchar NULL,
	expiry integer NULL,
	email varchar NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (email)
);

CREATE TABLE public.scheduled_emails (
	email_id varchar NOT NULL,
	scheduled_time integer NOT NULL,
	read_receipt boolean NULL,
	CONSTRAINT scheduled_emails_pk PRIMARY KEY email
);