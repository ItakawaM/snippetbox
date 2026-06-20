# snippetbox

## What is this?

Custom implementation of Snippetbox, a web application designed by [Alex Edwards](https://github.com/alexedwards) in his book "Let's Go"

## Features

- **Snippet Management**:
Create, view and browse text snippets.

- **User Management**:
Full signup, login, and logout flow with bcrypt password hashing and password strength validation. 

- **Account Management**:
View account details and change password.

- **Security**:
Self-signed SSL Certificate, TLS 1.2/1.3 enforced, CSRF protection via `nosurf` and CSRF Tokens.

- **Testing**:
Unit tests, end-to-end tests and intergration tests via a real PostreSQL database with setup and teardown scripts.

## New Features

This sections consists of changes/improvements I made while reading the book.

- **PostgreSQL instead of MySQL**
- **Golang 1.22 `net/http` `ServeMux` instead of a separate router**
- **Fully implemented comments system**:
Authenticated users can post comments under snippets, everyone can view them.
- **Environment Variables via `.env` and Build via `Makefile`**

## Afterword

Had a lot of fun reading the book!
