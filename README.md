
# Organization Service

Organization service for create/update/delete organization for user. It allows to invite/remove other users, assign them role and update their roles.

# Features
- User can create organizations.
- Every organization has 4 different roles of Owner, Admin, Viewer, Editor.
- Owner and Admin can invite different users to join their organization.
- Owner and Admin can change role of organization's member and also can remove them from organization.
- User can accept/reject the invitation.
- Organization's member can leave the organization.
- Organization's owner can delete the organization through otp.
- Organization's owner can transfer the ownership.

# Tech Stack 
- GO 1.21
- CockroachDB
- Dbmate
- JWT (json web token)

## Run Locally

Prerequisites you need to set up on your local computer:

- [Golang](https://go.dev/doc/install)
- [Cockroach](https://www.cockroachlabs.com/docs/releases/)
- [Dbmate](https://github.com/amacneil/dbmate#installation)

1. Clone the project

```bash
  git clone https://github.com/KaranLathiya/organization.git
  cd organization
```

2. Copy the .env.example file to new .config/.env file and set env variables in .env:

```bash
  cp .env.example .config/.env
```

3. Create `.env` file in current directory and update below configurations:
   - Add Cockroach database URL in `DATABASE_URL` variable.
4. Run `dbmate migrate` to migrate database schema.
5. Run `go run cmd/main.go` to run the programme.

# Routing

## For organization 

To create new organization  --POST

    http://localhost:9000/organization/
    
To update organization details  --PUT

    http://localhost:9000/organization/
    
To delete organization by owner (get otp)  --POST

    http://localhost:9000/organization/otp

## For invitation

To invite users to join the organizaton --POST

    http://localhost:9000/organization/members/invitation/

To accept/reject invitation of organizaton --DELETE

    http://localhost:9000/organization/member/invitation/
    
To get list of all invitations --GET

    http://localhost:9000/organization/member/invitation/

## For update organization's member role 

To update organization's member role (by owner/admin) --PUT

    http://localhost:9000/organization/members/role/
    
To transfer organization's ownership (by owner) --PUT

    http://localhost:9000/organization/members/role/owner

## For remove member from organization

To remove member from organization --DELETE

    http://localhost:9000/organization/members

## For organization details of user

To get all organizations details of user --GET

    http://localhost:9000/member/organizations
    
To get organization details of user --GET

    http://localhost:9000/member/{organization}/organization

## Public apis 

To jwt for  organizations details of users --GET

    http://localhost:9000/internal/members/organizations/jwt
    
To get organizations details of users  --POST

    http://localhost:9000/internal/members/organizations
    
To delete organization  --POST

    http://localhost:9000/internal/organization/{organization}


