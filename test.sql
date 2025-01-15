CREATE TABLE public.users (
	access_token varchar NULL,
	token_type varchar NULL,
	refresh_token varchar NULL,
	expiry integer NULL,
	email varchar NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (email)
);

CREATE TABLE public.scheduled_emails (
	message_id varchar NOT NULL,
	scheduled_time integer NOT NULL,
	read_receipt boolean NULL,
	username varchar NOT NULL, -- this is the email of the user
	CONSTRAINT scheduled_emails_pk PRIMARY KEY message_id
);

CREATE TABLE public.tracked_emails (
	message_id varchar NOT NULL,
	thread_id varchar NULL,
	username varchar NOT NULL, -- this is the email of the user
	CONSTRAINT tracked_emails_pk PRIMARY KEY message_id
);

-- create table to track the views of an email
CREATE TABLE public.email_views (
	message_id varchar NOT NULL,
	viewed_time integer NOT NULL,
	username varchar NOT NULL, -- this is the email of the user
	CONSTRAINT email_views_pk PRIMARY KEY message_id
);