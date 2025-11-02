# Merlin

## Game State

### Status

#### Status Tree & Outline

[Status Tree](docs_assets/status-tree-2.png)

##### I. Status Checkpoints<br>
    A. **Idle - Can** signify various points within the actual game include Home Screen, End of Game Screen, Settings, etc.<br>
    B. **Selection** - When a Game Mode (Quick Play or Competitive) has been selected and Roles are being chosen<br>
    C. **In Queue** - Status of waiting to get into a match<br>
    D. **Map Voting** - When the user first enters a match and everyone is voting for a map<br>
    E. **Banning Phase** - (Competitive Game Mode Only) The phase in which players are choosing their preferred hero and voting on what heroes to ban<br>
##### II. Natural Path - Ideal status flow<br>
    A. **Idle-Selection** - Entering role selection<br>
    B. **Selection-In Queue** - Roles selected, waiting for a match<br>
    C. **In Queue-Map Voting** - Entered a match, voting for a map<br>
    D. **Map Voting-Banning Phase** - Map selected, banning phase commenced<br>
    E. **Banning Phase-Idle** - No in game tooling so after banning phase we return to idle. Ideal situation is screen monitoring is turned off at this point<br>
##### III. Diversions - Not connected to each necessarily like "Natural Path", exceptions to ideal status flow<br>
    A. **Idle-In Queue** - Most commonly from the of end of match requeue selection<br>
    B. **In Queue-Selection** - Most commonly from cancelling of queue or changing of roles<br>
    C. **Selection-Idle** - Leaves role and queue selection<br>
##### IV. A Note on Entry Paths<br>
    - Every Status Checkpoint will have a way for the user to enter that status right away. I do not know when a user will start monitoring their screen
    so it is important to keep the entry point intuitive, because data can provided from them if the time has passed to capture it but what is most important
    is I am not breaking my system and I am creating a seamless UX.<br>