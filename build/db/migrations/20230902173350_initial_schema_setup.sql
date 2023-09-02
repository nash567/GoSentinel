-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.applications (
	id text NOT NULL,
	name text NOT NULL,
	email text NOT NULL,
	status text NOT NULL,
	is_verified bool NULL DEFAULT false,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	password text NULL,
	CONSTRAINT applications_email_key UNIQUE (email),
	CONSTRAINT applications_pkey PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS public.application_identities (
	id text NOT NULL,
	application_id text NOT NULL,
	secret varchar(255) NOT NULL,
	secret_viewed bool NULL DEFAULT false,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	password varchar NULL,
	CONSTRAINT application_identities_pkey PRIMARY KEY (id),
	CONSTRAINT application_identities_application_id_fkey FOREIGN KEY (application_id) REFERENCES public.applications(id)
);

-- +goose StatementEnd
