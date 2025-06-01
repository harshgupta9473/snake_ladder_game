# Snake And Ladder Game 
Built in golang following  clean layerd architecture.

This Game is multiplayer game (currently two player) , and it is played in realtime using websockets and packet based backend. Routes are decided based on packet type.
User can request to join a game and there is matchmaking service which joins two players who have opted for the same game type and starts the game, it is turn based game and the changes are broadcasted for every player of that game.
There is randomly generatd board with snakes and ladders which follows normal snake and ladder concept of conventional snake and ladder game.
If user leaves the game before it is finished that is before someone wins, the game will wait for atleast 30 seconds befoe ending and drawing the game,if user make connection to the backend he/she will be  joined to again the same game.
But if user doesn't tries making the reconnection in that 30 second interval, the game will end and the result will be draw.

This game is built in Golang following the clean layered architecture with different layers such as domain layer, repository layer,service layer and transport layer.
Service layer consists of three services: 
                                           UserService: Main work is to make websocket connection and save the conneciton,uses user repository and  transport layer for the connection management, and sending and recieving message to user.
                                           GameService: Main work is to Create and Start the game for the user, uses Game Repository for game creation and game board and User service for broadcasting game status to the players of that game.
                                           MatchMaking Service: Main work is to Match Make one user with other and start the game.
