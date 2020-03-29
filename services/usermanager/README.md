User Manager
=====

User manager is a service that keeps track of all player connection states. If a user disconnects, user manager will handle updating Redis to synchronize the user states globally. 

Calls that usermanager handles:


        
    disconnect <session id>
    
        Handles removing a user from the logged in sessions
       
        Will also notify the character state manager.
        
    connect <account> <character> <session id>
    
        Handles adding a user to the logged in sessions
        
        Will also notify the character state manager. 
        
    who <filter> 
    
        Will give a list of users that match the filter.
        
    account <account>
    
        Will check if a target user (account name) is logged in
        as well as details about the session.
        
    status <id> 
        
        Status of a given session ID 
    