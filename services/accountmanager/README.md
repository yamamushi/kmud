Account Manager
====

Handles account management and authorization.

Will also return account information. 

## API

    /auth
    
        Request: 
            Secret: (string) Secret shared token used by frontend service for Auth.
            Username: (string) Account Username
            HashedPass: (string) Sha256 hashed PW 
            
        Response:
            AuthToken: (string) Account Auth Token
            Error: Error status
    
    
    /accountinfo
    
        Request:
            Secret: (string) Secret shared token used by frontend service for Auth.
            AuthToken: (string) User Account Auth Token
            Field: (string) Target field(s) - Accepts all, email, permissions, characters, locked  
            
        Response: 
            Account: An account object with requested field(s)
            Error: Error status 
            
            
    /register
    
        Request:
            Secret: (string) Secret shared token used by frontend service for Auth.
            Username: (string) Account Username
            HashedPass: (string) Sha256 hashed PW 
            Email: (string) Account Email Address
            
        Response:
            Error: Error status (empty on success)
            
## Examples

### Register an Account

    curl -XPOST -d'{"secret":"secret","username":"accountusername","hashedpass":"hashedpassword","email":"account@email.com"}' localhost:4242/register
    
Example Output

    {"error":""}
<sub>Note: Output will be empty if registration was successful

Example Errors 

    {"error":"account with username accountusername already exists"}
    
    {"error":"account with email account@email.com already exists"}

    {"error":"unauthorized request"}

### Retrieve Auth Token 

    curl -XPOST -d'{"secret":"secret","username":"accountusername","hashedpass":"hashedpass"}' localhost:4242/auth
    
Example Output

    {"authtoken":"accountusername:EHfo458Yd7YvQFXWWmHkU1dXwcumbXg9oakfFQuXw"}

Example Errors

    {"authtoken":"","error":"invalid password"}
    
    {"authtoken":"","error":"account not found"}
    
    {"error":"unauthorized request"}
    
    
### Retrieve Account Info

Filter for all fields

    curl -XPOST -d'{"secret":"secret","token":"accountusername:H5rHuz382PfIVfLCt4EuKsJRohyrK5SuiyqyTErEo","field":"all"}' localhost:4242/accountinfo
    
Output
    
    {"account":{"username":"accountusername","email":"account@email.com","permissions":["user"],"locked":"false"}}
    
Filter for email

    curl -XPOST -d'{"secret":"secret","token":"accountusername:H5rHuz382PfIVfLCt4EuKsJRohyrK5SuiyqyTErEo","field":"email"}' localhost:4242/accountinfo
    
Output
    
    {"account":{"username":"accountusername","email":"account@email.com"}}
    
Example Errors

    {"account":{},"error":"unauthorized request"}

    {"account":{},"error":"invalid token format"}
    
    {"account":{},"error":"unrecognized fields: foobar"}

    