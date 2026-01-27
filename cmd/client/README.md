

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
│    └─ page_name
│         ├── view_function.go
│         └── http_request_funciton.go

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

