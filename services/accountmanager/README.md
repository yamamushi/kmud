Account Manager
====

Handles account management and authorization.

Will also return account information. 

API
====

    /auth
    
        Request: 
            Secret: (string) Secret shared token used by frontend service for Auth.
            Username: (string) Account Username
            HashedPassword: (string) Sha256 hashed PW 
            
        Response:
            AuthToken: (string) Account Auth Token
            Error: Error status
    
    
    /accountinfo
    
        Request:
            Secret: (string) Secret shared token used by frontend service for Auth.
            AuthToken: User Account Auth Token
            Account: Target account 
            Field: Target field
            
        Response: 
            Account: An account object with requested field(s)
            Error: Error status 