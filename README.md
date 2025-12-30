

## **Client architecture:**
* TUI + BubbleTea = state machine + message passing
* For each page - create new model
* Navigation - separate logic layer


## **Project structure:**
<pre>
tui/
├── main.go
├── app/
│    ├── model.go
│    ├── update.go
│    ├── view.go
│    └── messages.go
│
├── nav/
│    └── navigator.go
│
├── pages/
│    ├── page.go
│    ├── auth/
│    │    └── page.go          // 0 Login/Registration page
│    ├── main/
│    │    └── page.go          // 1 page where will be 2 actions: Load/Upload
│    ├── upload/
│    │    ├── page.go          // 2 list with possible types of data
│    │    ├── card/
│    │    │    └── page.go     // 2.1 page for card upload
│    │    ├── text/
│    │    │    └── page.go     // 2.2 page for simple text upload
│    │    └── file/
│    │         └── page.go     // 2.3 page for byte data upload
│    └── load/
│         ├── page.go          // 3 list with possible types of data
│         ├── card/
│         │    ├── page.go     // 3.1 page for list of cards objects
│         │    └── object/    
│         │          └── page.go     // 3.10 page for single object from list, where all fields are editable
│         ├── text/
│         │    └── page.go     // 3.2 page for list of simple text objects
│         │    └── object/    
│         │          └── page.go     // 3.20 page for single object from list, where all fields are editable
│         └── file/
│              └── page.go     // 3.3 page for list of byte bata objects
│              └── object/    
│                    └── page.go     // 3.30 page for single object from list, where all fields are editable
│
├── shared/
├── state.go
└── components/
</pre>

## **Navigation between pages**
_Navigation is implemented using a stack-based approach, allowing
natural support for nested pages and backward navigation._

##### Navigator
_Stack-based page manager_

* Push(page) — pushes a new page onto the stack and makes it active
* Pop() — removes the current page and returns to the previous one
* Current() — returns the currently active page (top of the stack)

Pages never interact with the navigator directly.
All navigation is performed via messages.k

##### Root model/view
_Central application coordinator_

* owns the Navigator
* stores global application state
* processes navigation messages
* delegates all other messages to the active page



The root view always renders the current page provided by the navigator

