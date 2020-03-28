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
            Error: (string) Error status
    
    
    /accountinfo
    
        Request:
            Secret: (string) Secret shared token used by frontend service for Auth.
            AuthToken: (string) User Account Auth Token
            Field: (string) Target field(s) - Accepts all, email, permissions, characters, locked  
            
        Response: 
            Account: (types.Account) An account object with requested field(s)
            Error: (string) Error status 
            
            
    /register
    
        Request:
            Secret: (string) Secret shared token used by frontend service for Auth.
            Username: (string) Account Username
            HashedPass: (string) Sha256 hashed PW 
            Email: (string) Account Email Address
            
        Response:
            Error: (string) Error status (empty on success)
            
    
    /search
    
        Request: Secret: (string) Secret shared token used by frontend service for Auth.
        AuthToken: (string) User Account Auth Token
        Account: (types.Account) Formatted account object to filter by
        
        Response:
            Accounts: ([]types.Account)   
            Error: (string) Error status (empty on success)  
    
    /modify
    
        Request: Secret: (string) Secret shared token used by frontend service for Auth.
        AuthToken: (string) User Account Auth Token
        Account: (types.Account) Formatted account object to modify (email or username)
            Note: You should stick to one field modification per use.
        
        Response:
            Account: (types.Account) Modified Account Record   
            Error: (string) Error status (empty on success) 
            
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
    
    {
    	"account": {
    		"username": "accountusername",
    		"email": "account@email.com",
    		"permissions": ["user"],
    		"locked": "false"
    	}
    }
    
Filter for email

    curl -XPOST -d'{"secret":"secret","token":"accountusername:H5rHuz382PfIVfLCt4EuKsJRohyrK5SuiyqyTErEo","field":"email"}' localhost:4242/accountinfo
    
Output
    
    {
    	"account": {
    		"username": "accountusername",
    		"email": "account@email.com"
    	}
    }
    
Example Errors

    {"account":{},"error":"unauthorized request"}

    {"account":{},"error":"invalid token format"}
    
    {"account":{},"error":"unrecognized fields: foobar"}

    
### Search For Account

    curl -XPOST -d'{"secret":"secret","token":"accountusername:H5rHuz382PfIVfLCt4EuKsJRohyrK5SuiyqyTErEo","account":{"permissions":["user"]}}' localhost:4242/search
    
Output

    {
    	"accounts": [{
    		"username": "me",
    		"email": "me@email.com",
    		"hashedpass": "74657374",
    		"permissions": ["user"],
    		"groups": ["default"],
    		"characters": ["mycharacter1"],
    		"locked": "false",
    	}, {
    		"username": "you",
    		"email": "you@mail.com",
    		"hashedpass": "74657374",
    		"permissions": ["user"],
            "groups": ["default"],
            "locked": "false",
    	}]
    }

Example Errors

    {"accounts":[],"error":"unauthorized request"}
    
    {"account":{},"error":"invalid token format"}
    
    
### Modify Account

    curl -XPOST -d'{"secret":"secret420","token":"yamamushi2001:gSvJnwml38wliLKmspFOh2moNEewAiMRvRgc3CW6A","account":{"username":"accountname","email":"newemail@email.com","hashedpass":"74657374"}}' localhost:4242/modify
    
Output

    {
    	"account": {
    		"username": "accountname",
    		"email": "newemail@email.com",
    		"hashedpass": "",
    		"groups": ["default", "moderators"],
    		"permissions": ["user"],
    		"locked": "false",
    		"token": ""
    	}
    }

Example Errors

    {"accounts":[],"error":"unauthorized request"}
    
    {"account":{},"error":"invalid token format"}    