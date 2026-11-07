# Todo: Taking a Hiatus
Up to this point the project still needs the following:
 - 1 Sesssion Rotation: should be atomic (rotate insde the db transaction)
 - 2 Session Expiry: scheduled worker to remove expired sessions
 - 3 JWT revocation: add token/session revocation checks for immediate revocation
 - 4 Organisation isolation: Never assume that a valid permission in one organisation grants access to another organisation's resources. Check the organisation on every relevant operation.
 - 5 JSON: Store JSON-compatible values in the session data. Try avoid storing go structs, 
 database models or any sensitive data like payment information