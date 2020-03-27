User Manager
=====

User manager is a service that keeps track of all player connection states. If a user disconnects, user manager will handle updating Redis to synchronize the user states globally. 

Calls that usermanager handles:

    globalwho 
    
        Will give a global list of who is logged in.
        
    disconnect
    
        Handles removing a user from the logged in sessions
       
        Will also notify the character state manager.
        
    connect 
    
        Handles adding a user to the logged in sessions
        
        Will also notify the character state manager. 
        
    accountstatus <name>
    
        Will check if a target user (account name) is logged in
        as well as details about the session.
        
    
    