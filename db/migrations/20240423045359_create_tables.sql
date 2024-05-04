-- migrate:up

CREATE TYPE IF NOT EXISTS public.privacy AS ENUM (
	'public',
	'private');

CREATE TYPE IF NOT EXISTS public."role" AS ENUM (
	'owner',
	'admin',
	'editor',
	'viewer');

CREATE TABLE IF NOT EXISTS public.organization (
	id VARCHAR(20) NOT NULL DEFAULT unique_rowid(),
	name VARCHAR(50) NOT NULL,
	owner_id VARCHAR(20) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp():::TIMESTAMP,
	privacy organization.public.privacy NOT NULL,
	updated_by VARCHAR(20) NULL,
	updated_at TIMESTAMP NULL,
	CONSTRAINT pk_organization_id PRIMARY KEY (id ASC)
);

CREATE TABLE IF NOT EXISTS public.invitation (
	id VARCHAR(20) NOT NULL DEFAULT unique_rowid(),
	invitee VARCHAR(20) NOT NULL,
	"role" organization.public."role" NOT NULL,
	organization_id VARCHAR(20) NOT NULL,
	invited_at TIMESTAMP NOT NULL DEFAULT current_timestamp():::TIMESTAMP,
	invited_by VARCHAR(20) NOT NULL,
	CONSTRAINT pk_invitation_id PRIMARY KEY (id ASC),
	CONSTRAINT fk_invitation_organization_id FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
	UNIQUE INDEX uc_invitation_invitee_and_organization_id (invitee ASC, organization_id ASC)
);

CREATE TABLE IF NOT EXISTS public.member (
	id VARCHAR(20) NOT NULL DEFAULT unique_rowid(),
	user_id VARCHAR(20) NOT NULL,
	"role" organization.public."role" NOT NULL,
	organization_id VARCHAR(20) NOT NULL,
	joined_at TIMESTAMP NOT NULL DEFAULT current_timestamp():::TIMESTAMP,
	updated_by VARCHAR(20) NULL,
	updated_at TIMESTAMP NULL,
	invited_by VARCHAR(20) NULL,
	CONSTRAINT pk_member_id PRIMARY KEY (id ASC),
	CONSTRAINT fk_member_organization_id FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
	UNIQUE INDEX uc_member_member_id_and_organization_id (user_id ASC, organization_id ASC)
);

-- migrate:down

DROP TABLE IF EXISTS public.invitation;

DROP TABLE IF EXISTS public.member;

DROP TABLE IF EXISTS public.organization;

DROP TYPE IF EXISTS public.privacy;

DROP TYPE IF EXISTS public."role";
